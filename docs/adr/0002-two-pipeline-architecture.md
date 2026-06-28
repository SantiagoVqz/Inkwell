# Two-pipeline architecture with SQLite as the seam

Inkwell is structured as two independent pipelines sharing one SQLite store, rather than one monolithic pipeline. The **ingest pipeline** (per-Entry, frequent, no LLM) handles fetch → normalize → write note → embed → store. The **synthesis pipeline** (per-window, infrequent, LLM-only) handles query vectors → cluster → synthesize → write Storyline note. The two never call each other directly; SQLite is the queue and source of truth between them.

This is the consequence of a deliberate AI-stack rethink: K-means / per-entry LLM categorization was rejected in favor of embeddings-for-cheap-deterministic-work + LLM-for-synthesis-only. Embeddings handle dedup, retrieval, and similarity (cheap, local-able, deterministic); clustering discovers themes; the LLM runs only once per cluster, not per Entry — roughly 100× fewer calls than the naive design.

**Considered and rejected**: single pipeline with per-Entry LLM categorization (expensive, slow, requires network on every ingest); K-means clustering specifically (requires K up front, no noise category, forces spherical clusters in embedding space — HDBSCAN or agglomerative-threshold are correct choices when we get to v3).
