# Opportunity Template Guide

## Purpose

Every opportunity in Banter Engine MUST follow the same structure.

## Template

```yaml
id: PredictionMassacre

name: Prediction Massacre

category: Crowd

severity: 5

description: >
  Large majority of players predicted the same outcome
  and that outcome failed.

requiredData:
  - forecasts
  - matchResults

trigger:
  majorityPercentage: 70
  majorityWrong: true

banterAngles:
  - crowd_humiliation
  - sheep_mentality
  - collective_failure

relatedOpportunities:
  - CrowdTrap
  - MajorityCollapse

tags:
  - crowd
  - upset
```

## Field Definitions

### id
Unique technical identifier.

### name
Human readable name.

### category
One of:
- Ranking
- Predictions
- Crowd
- MatchEvents
- Narratives

### severity
Scale from 1 to 5.

### description
Business description of the opportunity.

### requiredData
Data required to detect the opportunity.

### trigger
Detection conditions.

### banterAngles
Possible banter directions.

### relatedOpportunities
Related opportunities.

### tags
Search and classification tags.
