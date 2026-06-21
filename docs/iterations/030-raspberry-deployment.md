# Iteration 030 - Raspberry Deployment

## Goal

Deploy Banter Engine on Raspberry Pi and validate long-running execution.

## Tasks

- Create deployment guide
- Create systemd service
- Create .env.example
- Validate startup after reboot
- Validate graceful shutdown
- Validate automatic restart
- Validate snapshot persistence
- Validate network recovery
- Add operational runbook

## Deliverables

docs/raspberry-deployment.md
deploy/banter-engine.service
.env.example

## Definition of Done

- Service starts automatically
- Service survives reboot
- Service recovers from failures
- Messages reach Discord
- No manual intervention required
