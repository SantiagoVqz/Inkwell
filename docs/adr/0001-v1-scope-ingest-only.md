# v1 scope is ingest-only — no embeddings, no LLM

Inkwell is positioned as an "AI tool" but v1 ships with zero ML and zero LLM calls. We chose this because (a) ingest reliability is the only thing that proves the Go rewrite was worth doing, and the rest of the value can ride on top later; (b) it teaches the most Go per unit of code (concurrency, context, errors, sqlc, cobra) without entangling ML learning; (c) the existing Python `content-pipeline` keeps doing synthesis in parallel until Inkwell catches up, so no capability is lost.

Embeddings move to v2; clustering and LLM synthesis move to v3. The two-pipeline architecture (see ADR-0002) is designed so they slot in additively without rework.

**Considered and rejected**: shipping embeddings in v1 (pulls in unfamiliar work and an external provider dependency before the fetch path is proven); shipping the whole vision in v1 (long path to first usable release; bad fit for a learning project).
