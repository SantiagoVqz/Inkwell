## Language

**Feed**:
A subscribed RSS/Atom source identified by a URL. Inkwell fetches each Feed on a schedule and produces Entries from its items.
_Avoid_: Source, channel, subscription.

**Entry**:
One item from one Feed, captured at fetch time. Immutable once written. Each Entry corresponds to exactly one Obsidian note in the vault.
_Avoid_: Article, post, item, note (the markdown file is the *note*, the database row is the *Entry*).

**Story**:
A developing narrative thread that accumulates Entries over its lifetime (e.g. "Anthropic fundraise", "OpenAI lawsuit"). Has a status (open/closed). Many Entries may attach to one Story; one Entry may attach to zero or one Story.
_Avoid_: Topic, thread, theme, cluster (a *cluster* is a v3 synthesis-time concept and is different — see below).

**Storyline**:
A weekly (or window-based) synthesis narrative produced by the v3 synthesis pipeline. References multiple Stories and standalone Entries. Lives as a single markdown note in the vault. One Storyline per synthesis run.
_Avoid_: Digest, weekly summary, report.

**Cluster**:
A v3 synthesis-time grouping of Entries discovered by unsupervised clustering (HDBSCAN / agglomerative) over a time window. Distinct from a Story: Stories are persistent and human-named; Clusters are ephemeral and algorithm-named. Clusters may *promote* into Stories, but they are not the same.
_Avoid_: Topic, group, theme.

**Vault**:
The Obsidian directory tree on disk where Inkwell writes notes. Configured via `vault_path`. Inkwell only writes inside a configured subfolder; it never touches notes outside it.
_Avoid_: Repository, library, store.

**Ingest pipeline**:
The fetch → normalize → write-note → (v2: embed) → store path. Runs frequently (cron, e.g. every 6h). Per-Entry. No LLM, ever.
_Avoid_: Pipeline (ambiguous), fetcher (only one stage).

**Synthesis pipeline** *(v3)*:
The query-vectors → cluster → LLM-summarize → write-Storyline path. Runs infrequently (cron, e.g. weekly). Per-window. LLM lives here and only here.
_Avoid_: Enrichment (too vague), summarization (only one stage).

**Attachment**:
The link between an Entry and a Story. Has a `source` field recording HOW the attachment was made: `frontmatter`, `cli`, `embedding` (v2), or `llm` (v3).
_Avoid_: Linking, tagging (tags are a separate concept reserved for v2+).

## Cardinality summary

```
Feed 1 ─── N Entry
Entry 0..1 ── 1 Story     (Attachment)
Story 1 ─── N Entry
Storyline N ── N Story    (only in v3)
Storyline N ── N Entry    (only in v3; standalone Entries can be referenced)
Cluster N ── N Entry      (ephemeral, only during a synthesis run)
```

## Example dialogue

> **Dev:** Should the ingest pipeline call the LLM if a new Entry might belong to an open Story?
>
> **Domain:** No. Ingest never calls the LLM. In v1 you attach manually via frontmatter or CLI. In v2 the ingest pipeline runs an embedding similarity check against open Stories. The LLM only ever runs inside synthesis.
>
> **Dev:** So if I want to find related stuff at ingest time without the LLM, I use embeddings?
>
> **Domain:** In v2, yes. v1 has no automatic relation-finding at all — Entries are independent rows until you manually attach them to a Story.
>
> **Dev:** And a Storyline is just a fancy Story?
>
> **Domain:** No, they're different. A Story is something you (or eventually the embedder) declare exists and accumulate Entries into over weeks. A Storyline is a one-shot weekly synthesis output — it references Stories and ungrouped Entries to tell "what happened this week."
