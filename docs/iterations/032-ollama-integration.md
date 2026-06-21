# Iteration 031 - Ollama Integration

## Goal

Introduce local LLM generation using Ollama.

## Tasks

- Install Ollama on Raspberry Pi
- Benchmark candidate models
- Create Prompt Builder
- Create Ollama client
- Generate structured JSON responses
- Integrate catalog opportunities
- Integrate narratives
- Integrate archetypes
- Fallback to deterministic banter

## Candidate Models

- qwen3:4b
- qwen3:8b
- gemma3:4b

## Expected Flow

Opportunity
+
Context
+
Narrative
+
Archetype
↓
Prompt Builder
↓
Ollama
↓
Structured JSON
↓
Discord

## Definition of Done

- Local generation works
- Structured JSON validated
- Discord receives AI-generated banter
- Fallback mode available
