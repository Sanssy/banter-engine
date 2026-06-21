package engine

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/Sanssy/banter-engine/internal/banter"
	"github.com/Sanssy/banter-engine/internal/catalog"
	"github.com/Sanssy/banter-engine/internal/config"
	"github.com/Sanssy/banter-engine/internal/discord"
	"github.com/Sanssy/banter-engine/internal/forecasts"
	"github.com/Sanssy/banter-engine/internal/logging"
	matchmodel "github.com/Sanssy/banter-engine/internal/matches"
	"github.com/Sanssy/banter-engine/internal/mpp"
	"github.com/Sanssy/banter-engine/internal/opportunities"
	"github.com/Sanssy/banter-engine/internal/rivalries"
	"github.com/Sanssy/banter-engine/internal/snapshot"
)

const opportunityCatalogPath = "resources/opportunities.json"

type Engine struct {
	config  config.Config
	mpp     *mpp.Client
	discord *discord.Client
	catalog *catalog.Catalog
	logger  *logging.Logger
	output  io.Writer
}

func New(cfg config.Config, output io.Writer) (*Engine, error) {
	opportunityCatalog, err := catalog.LoadCatalog(opportunityCatalogPath)
	if err != nil {
		return nil, err
	}
	for _, opportunityType := range opportunities.RegisteredTypes() {
		if _, found := opportunityCatalog.FindByID(opportunityType); !found {
			return nil, fmt.Errorf("registered opportunity %q is missing from catalog", opportunityType)
		}
	}

	var discordClient *discord.Client
	if !cfg.DryRun {
		discordClient = discord.NewClient(cfg.DiscordWebhookURL)
	}
	return &Engine{
		config:  cfg,
		mpp:     mpp.NewClient(cfg.MPPToken),
		discord: discordClient,
		catalog: opportunityCatalog,
		logger:  logging.New(output, "engine"),
		output:  output,
	}, nil
}

func (e *Engine) Run(ctx context.Context) error {
	e.logger.Info("scheduler started with poll interval %s", e.config.PollInterval)
	defer e.logger.Info("scheduler stopped")

	ticker := time.NewTicker(e.config.PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		if err := e.runOnce(); err != nil {
			if ctx.Err() != nil {
				return nil
			}
			e.logger.Error("run failed: %v", err)
		}

		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
		}
	}
}

func (e *Engine) runOnce() error {
	previousStandings, err := snapshot.LoadStandings(e.snapshotPath("standings.json"))
	if err != nil {
		return err
	}
	previousMatches, err := snapshot.LoadMatches(e.snapshotPath("matches.json"))
	if err != nil {
		return err
	}
	rivalryState, err := snapshot.LoadRivalries(e.snapshotPath("rivalries.json"))
	if err != nil {
		return err
	}
	previousForecasts, err := snapshot.LoadForecasts(e.snapshotPath("forecasts.json"))
	if err != nil {
		return err
	}

	standings, err := e.mpp.GetStandings(e.config.ChallengeID)
	if err != nil {
		return err
	}
	matches, err := e.mpp.GetMatches()
	if err != nil {
		return err
	}

	previousMatchesByID := make(map[string]matchmodel.Match, len(previousMatches))
	for _, match := range previousMatches {
		previousMatchesByID[match.MatchID] = match
	}
	for i := range matches {
		isLive := matches[i].Status != "" && matches[i].Status != "preMatch" && matches[i].Status != "fullTime"
		previousMatch, existed := previousMatchesByID[matches[i].MatchID]
		justEnded := existed && matches[i].Status == "fullTime" && previousMatch.Status != "fullTime"
		if !isLive && !justEnded {
			continue
		}
		events, err := e.mpp.GetMatchEvents(matches[i].MatchID)
		if err != nil {
			return err
		}
		matches[i].Events = events
	}

	var forecastHistory []forecasts.Forecast
	var allForecasts []forecasts.Forecast
	var detected []opportunities.Opportunity
	updatedRivalries, rivalryOpportunities := rivalries.Update(standings, rivalryState)
	detected = append(detected, rivalryOpportunities...)
	detected = append(detected, opportunities.DetectLiveUpdates(previousMatches, matches)...)
	for _, match := range matches {
		matchForecasts, err := e.mpp.GetForecasts(e.config.ChallengeID, match)
		if err != nil {
			return err
		}
		allForecasts = append(allForecasts, matchForecasts...)
		if match.Status == "fullTime" {
			forecastHistory = append(forecastHistory, matchForecasts...)
		}
		detected = append(detected, opportunities.DetectSurprises(match, matchForecasts)...)
		if previousMatch, ok := previousMatchesByID[match.MatchID]; ok {
			detected = append(detected, opportunities.DetectHeartbreaks(previousMatch, match, matchForecasts)...)
			detected = append(detected, opportunities.DetectLatePointImpacts(previousMatch, match, previousForecasts, matchForecasts)...)
		}
	}
	detected = append(detected, opportunities.DetectPointImpacts(previousForecasts, allForecasts)...)
	detected = append(detected, opportunities.DetectStreaks(forecastHistory)...)
	detected = append(detected, opportunities.Detect(previousStandings, standings)...)

	for _, opportunity := range detected {
		definition, found := e.catalog.FindByID(opportunity.Type)
		if !found {
			return fmt.Errorf("unknown opportunity %q", opportunity.Type)
		}
		message := banter.GenerateWithDefinition(opportunity, definition)
		if err := e.publish(message); err != nil {
			return err
		}
	}

	if err := snapshot.SaveStandings(e.snapshotPath("standings.json"), standings); err != nil {
		return err
	}
	if err := snapshot.SaveMatches(e.snapshotPath("matches.json"), matches); err != nil {
		return err
	}
	if err := snapshot.SaveRivalries(e.snapshotPath("rivalries.json"), updatedRivalries); err != nil {
		return err
	}
	if err := snapshot.SaveForecasts(e.snapshotPath("forecasts.json"), allForecasts); err != nil {
		return err
	}
	e.logger.Info("run completed with %d opportunities", len(detected))
	return nil
}

func (e *Engine) publish(message string) error {
	if e.config.DryRun {
		_, err := fmt.Fprintln(e.output, message)
		return err
	}
	return e.discord.Send(message)
}

func (e *Engine) snapshotPath(name string) string {
	return filepath.Join(e.config.SnapshotDir, name)
}
