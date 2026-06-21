# LLM Layer Review

Iteration 029 found no active LLM provider or Ollama adapter. `internal/banter/llm.go` was instantiated only with a `nil` provider, so the provider branch could never run. `internal/context` and `internal/contextbuilder` existed solely to build input for that unused branch. The referenced `internal/banter/provider.go` file did not exist.

The unused LLM abstraction, its tests, and both context packages were removed. The engine now calls the deterministic banter generator directly.

Future Ollama work should introduce a provider interface together with its first real adapter and runtime configuration. The catalog already contains the metadata needed to build that integration without retaining inactive scaffolding.
