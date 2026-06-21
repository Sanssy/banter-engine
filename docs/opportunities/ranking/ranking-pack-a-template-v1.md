# Pack A - Ranking Opportunities

---
id: RankingOvertake

name: Ranking Overtake

category: Ranking

severity: 2

description: >
  A player overtakes another player.

requiredData:
  - standings

trigger:
  rankImprovement: 1
overtakenPlayers: 1

banterAngles:
  - superiority
  - overtake
  - humiliation

relatedOpportunities:
  - DoubleOvertake
  - BiggestClimber

tags:
  - ranking
  - movement

---
id: DoubleOvertake

name: Double Overtake

category: Ranking

severity: 3

description: >
  A player overtakes two players.

requiredData:
  - standings

trigger:
  rankImprovement: 2
overtakenPlayers: 2

banterAngles:
  - momentum
  - dominance

relatedOpportunities:
  - TripleOvertake
  - BiggestClimber

tags:
  - ranking
  - movement

---
id: TripleOvertake

name: Triple Overtake

category: Ranking

severity: 4

description: >
  A player overtakes three or more players.

requiredData:
  - standings

trigger:
  rankImprovement: 3
overtakenPlayers: 3

banterAngles:
  - shock
  - dominance

relatedOpportunities:
  - BiggestClimber
  - ComebackSeason

tags:
  - ranking
  - surge

---
id: EnteredTop3

name: Entered Top 3

category: Ranking

severity: 3

description: >
  A player enters the podium.

requiredData:
  - standings

trigger:
  newRankMax: 3

banterAngles:
  - glory
  - arrival

relatedOpportunities:
  - NewLeader
  - PodiumGuardian

tags:
  - ranking
  - podium

---
id: ExitedTop3

name: Exited Top 3

category: Ranking

severity: 3

description: >
  A player leaves the podium.

requiredData:
  - standings

trigger:
  oldRankMax: 3

banterAngles:
  - fall
  - pressure

relatedOpportunities:
  - FreeFall
  - PodiumGuardian

tags:
  - ranking
  - podium

---
id: EnteredTop10

name: Entered Top 10

category: Ranking

severity: 2

description: >
  A player enters top ten.

requiredData:
  - standings

trigger:
  newRankMax: 10

banterAngles:
  - recognition

relatedOpportunities:
  - RankingOvertake

tags:
  - ranking
  - top10

---
id: ExitedTop10

name: Exited Top 10

category: Ranking

severity: 2

description: >
  A player exits top ten.

requiredData:
  - standings

trigger:
  oldRankMax: 10

banterAngles:
  - decline

relatedOpportunities:
  - FreeFall

tags:
  - ranking
  - top10

---
id: NewLeader

name: New Leader

category: Ranking

severity: 5

description: >
  Leadership changes.

requiredData:
  - standings

trigger:
  rank: 1

banterAngles:
  - crowning
  - dominance

relatedOpportunities:
  - LeaderUnderPressure

tags:
  - ranking
  - leader

---
id: LeaderUnderPressure

name: Leader Under Pressure

category: Ranking

severity: 3

description: >
  Leader has tiny advantage.

requiredData:
  - standings

trigger:
  gapMax: 10

banterAngles:
  - stress
  - chase

relatedOpportunities:
  - NewLeader
  - PhotoFinish

tags:
  - leader
  - gap

---
id: RunawayLeader

name: Runaway Leader

category: Ranking

severity: 4

description: >
  Leader has huge advantage.

requiredData:
  - standings

trigger:
  gapMin: 100

banterAngles:
  - untouchable
  - dominance

relatedOpportunities:
  - NewLeader

tags:
  - leader
  - gap

---
id: PhotoFinish

name: Photo Finish

category: Ranking

severity: 3

description: >
  Very small gap between rivals.

requiredData:
  - standings

trigger:
  gapMax: 5

banterAngles:
  - tension
  - race

relatedOpportunities:
  - DeadHeat

tags:
  - ranking
  - close

---
id: DeadHeat

name: Dead Heat

category: Ranking

severity: 2

description: >
  Players have equal points.

requiredData:
  - standings

trigger:
  equalPoints: true

banterAngles:
  - stalemate

relatedOpportunities:
  - PhotoFinish

tags:
  - ranking
  - tie

---
id: LastPlaceLocked

name: Last Place Locked

category: Ranking

severity: 2

description: >
  Same player stays last for long period.

requiredData:
  - standings

trigger:
  weeksLastPlace: 4

banterAngles:
  - mockery
  - struggle

relatedOpportunities:
  - EscapedLastPlace

tags:
  - ranking
  - bottom

---
id: EscapedLastPlace

name: Escaped Last Place

category: Ranking

severity: 3

description: >
  Player leaves last place.

requiredData:
  - standings

trigger:
  leftLastPlace: true

banterAngles:
  - survival
  - relief

relatedOpportunities:
  - LastPlaceLocked

tags:
  - ranking
  - bottom

---
id: ComebackSeason

name: Comeback Season

category: Ranking

severity: 4

description: >
  Large positive ranking recovery.

requiredData:
  - standings

trigger:
  placesGained: 10

banterAngles:
  - redemption

relatedOpportunities:
  - BiggestClimber

tags:
  - ranking
  - comeback

---
id: FreeFall

name: Free Fall

category: Ranking

severity: 4

description: >
  Large ranking collapse.

requiredData:
  - standings

trigger:
  placesLost: 10

banterAngles:
  - collapse
  - panic

relatedOpportunities:
  - BiggestDropper

tags:
  - ranking
  - collapse

---
id: BiggestClimber

name: Biggest Climber

category: Ranking

severity: 3

description: >
  Best weekly progression.

requiredData:
  - standings

trigger:
  weeklyRankGainMax: true

banterAngles:
  - momentum

relatedOpportunities:
  - ComebackSeason

tags:
  - ranking
  - weekly

---
id: BiggestDropper

name: Biggest Dropper

category: Ranking

severity: 3

description: >
  Worst weekly regression.

requiredData:
  - standings

trigger:
  weeklyRankLossMax: true

banterAngles:
  - failure

relatedOpportunities:
  - FreeFall

tags:
  - ranking
  - weekly

---
id: PodiumGuardian

name: Podium Guardian

category: Ranking

severity: 3

description: >
  Long podium presence.

requiredData:
  - standings

trigger:
  weeksOnPodium: 4

banterAngles:
  - consistency

relatedOpportunities:
  - EnteredTop3

tags:
  - ranking
  - podium

---
id: OnePointBehind

name: One Point Behind

category: Ranking

severity: 2

description: >
  Only one point separates rivals.

requiredData:
  - standings

trigger:
  gapExact: 1

banterAngles:
  - pressure
  - chase

relatedOpportunities:
  - PhotoFinish

tags:
  - ranking
  - close
