# Pack B - Prediction Opportunities

---
id: PerfectPrediction

name: Perfect Prediction

category: Predictions

severity: 5

description: >
  Exact score prediction achieved.

requiredData:
  - forecasts
  - matchResults

trigger:
  detectionRule: PerfectPrediction

banterAngles:
  - predictions
  - confidence
  - accuracy

relatedOpportunities:
  - PredictionKing
  - PredictionFraud

tags:
  - predictions
  - forecasts

---
id: AlmostPerfectPrediction

name: Almost Perfect Prediction

category: Predictions

severity: 3

description: >
  One goal away from exact score.

requiredData:
  - forecasts
  - matchResults

trigger:
  detectionRule: AlmostPerfectPrediction

banterAngles:
  - predictions
  - confidence
  - accuracy

relatedOpportunities:
  - PredictionKing
  - PredictionFraud

tags:
  - predictions
  - forecasts

---
id: LuckyPrediction

name: Lucky Prediction

category: Predictions

severity: 2

description: >
  Correct outcome despite low confidence.

requiredData:
  - forecasts
  - matchResults

trigger:
  detectionRule: LuckyPrediction

banterAngles:
  - predictions
  - confidence
  - accuracy

relatedOpportunities:
  - PredictionKing
  - PredictionFraud

tags:
  - predictions
  - forecasts

---
id: DisasterPrediction

name: Disaster Prediction

category: Predictions

severity: 4

description: >
  Prediction completely misses reality.

requiredData:
  - forecasts
  - matchResults

trigger:
  detectionRule: DisasterPrediction

banterAngles:
  - predictions
  - confidence
  - accuracy

relatedOpportunities:
  - PredictionKing
  - PredictionFraud

tags:
  - predictions
  - forecasts

---
id: LastMinutePrediction

name: Last Minute Prediction

category: Predictions

severity: 2

description: >
  Forecast submitted just before kickoff.

requiredData:
  - forecasts
  - matchResults

trigger:
  detectionRule: LastMinutePrediction

banterAngles:
  - predictions
  - confidence
  - accuracy

relatedOpportunities:
  - PredictionKing
  - PredictionFraud

tags:
  - predictions
  - forecasts

---
id: WeekendMaster

name: Weekend Master

category: Predictions

severity: 4

description: >
  Outstanding prediction performance during a round.

requiredData:
  - forecasts
  - matchResults

trigger:
  detectionRule: WeekendMaster

banterAngles:
  - predictions
  - confidence
  - accuracy

relatedOpportunities:
  - PredictionKing
  - PredictionFraud

tags:
  - predictions
  - forecasts

---
id: WeekendDisaster

name: Weekend Disaster

category: Predictions

severity: 4

description: >
  Very poor prediction performance during a round.

requiredData:
  - forecasts
  - matchResults

trigger:
  detectionRule: WeekendDisaster

banterAngles:
  - predictions
  - confidence
  - accuracy

relatedOpportunities:
  - PredictionKing
  - PredictionFraud

tags:
  - predictions
  - forecasts

---
id: ScoreOracle

name: Score Oracle

category: Predictions

severity: 5

description: >
  Multiple exact scores achieved in same round.

requiredData:
  - forecasts
  - matchResults

trigger:
  detectionRule: ScoreOracle

banterAngles:
  - predictions
  - confidence
  - accuracy

relatedOpportunities:
  - PredictionKing
  - PredictionFraud

tags:
  - predictions
  - forecasts

---
id: DrawSpecialist

name: Draw Specialist

category: Predictions

severity: 3

description: >
  Repeated success predicting draws.

requiredData:
  - forecasts
  - matchResults

trigger:
  detectionRule: DrawSpecialist

banterAngles:
  - predictions
  - confidence
  - accuracy

relatedOpportunities:
  - PredictionKing
  - PredictionFraud

tags:
  - predictions
  - forecasts

---
id: FavoriteHunter

name: Favorite Hunter

category: Predictions

severity: 3

description: >
  Repeated success on favorites.

requiredData:
  - forecasts
  - matchResults

trigger:
  detectionRule: FavoriteHunter

banterAngles:
  - predictions
  - confidence
  - accuracy

relatedOpportunities:
  - PredictionKing
  - PredictionFraud

tags:
  - predictions
  - forecasts

---
id: UnderdogBeliever

name: Underdog Believer

category: Predictions

severity: 4

description: >
  Repeated success on underdogs.

requiredData:
  - forecasts
  - matchResults

trigger:
  detectionRule: UnderdogBeliever

banterAngles:
  - predictions
  - confidence
  - accuracy

relatedOpportunities:
  - PredictionKing
  - PredictionFraud

tags:
  - predictions
  - forecasts

---
id: RiskReward

name: Risk Reward

category: Predictions

severity: 3

description: >
  High-risk prediction pays off.

requiredData:
  - forecasts
  - matchResults

trigger:
  detectionRule: RiskReward

banterAngles:
  - predictions
  - confidence
  - accuracy

relatedOpportunities:
  - PredictionKing
  - PredictionFraud

tags:
  - predictions
  - forecasts

---
id: SafeBetFailure

name: Safe Bet Failure

category: Predictions

severity: 4

description: >
  Consensus safe prediction fails.

requiredData:
  - forecasts
  - matchResults

trigger:
  detectionRule: SafeBetFailure

banterAngles:
  - predictions
  - confidence
  - accuracy

relatedOpportunities:
  - PredictionKing
  - PredictionFraud

tags:
  - predictions
  - forecasts

---
id: SafeBetSuccess

name: Safe Bet Success

category: Predictions

severity: 2

description: >
  Consensus safe prediction succeeds.

requiredData:
  - forecasts
  - matchResults

trigger:
  detectionRule: SafeBetSuccess

banterAngles:
  - predictions
  - confidence
  - accuracy

relatedOpportunities:
  - PredictionKing
  - PredictionFraud

tags:
  - predictions
  - forecasts

---
id: ExactScoreMachine

name: Exact Score Machine

category: Predictions

severity: 5

description: >
  Many exact scores over time.

requiredData:
  - forecasts
  - matchResults

trigger:
  detectionRule: ExactScoreMachine

banterAngles:
  - predictions
  - confidence
  - accuracy

relatedOpportunities:
  - PredictionKing
  - PredictionFraud

tags:
  - predictions
  - forecasts

---
id: ExactScoreDrought

name: Exact Score Drought

category: Predictions

severity: 2

description: >
  Long period without exact scores.

requiredData:
  - forecasts
  - matchResults

trigger:
  detectionRule: ExactScoreDrought

banterAngles:
  - predictions
  - confidence
  - accuracy

relatedOpportunities:
  - PredictionKing
  - PredictionFraud

tags:
  - predictions
  - forecasts

---
id: PredictionKing

name: Prediction King

category: Predictions

severity: 5

description: >
  Best predictor in standings.

requiredData:
  - forecasts
  - matchResults

trigger:
  detectionRule: PredictionKing

banterAngles:
  - predictions
  - confidence
  - accuracy

relatedOpportunities:
  - PredictionKing
  - PredictionFraud

tags:
  - predictions
  - forecasts

---
id: PredictionFraud

name: Prediction Fraud

category: Predictions

severity: 4

description: >
  Persistently poor predictor.

requiredData:
  - forecasts
  - matchResults

trigger:
  detectionRule: PredictionFraud

banterAngles:
  - predictions
  - confidence
  - accuracy

relatedOpportunities:
  - PredictionKing
  - PredictionFraud

tags:
  - predictions
  - forecasts

---
id: AllOrNothing

name: All Or Nothing

category: Predictions

severity: 3

description: >
  Extreme variance between rounds.

requiredData:
  - forecasts
  - matchResults

trigger:
  detectionRule: AllOrNothing

banterAngles:
  - predictions
  - confidence
  - accuracy

relatedOpportunities:
  - PredictionKing
  - PredictionFraud

tags:
  - predictions
  - forecasts

---
id: PredictionSniper

name: Prediction Sniper

category: Predictions

severity: 4

description: >
  Rare predictions with high precision.

requiredData:
  - forecasts
  - matchResults

trigger:
  detectionRule: PredictionSniper

banterAngles:
  - predictions
  - confidence
  - accuracy

relatedOpportunities:
  - PredictionKing
  - PredictionFraud

tags:
  - predictions
  - forecasts
