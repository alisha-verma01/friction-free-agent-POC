FrictionFreeAgent — Agentic AI to Reduce Operational Friction in the Provider Portal

An internship project that builds agentic AI workflows to keep provider operations moving—even when core services are slow or down. The system automates CPT prior-authorization checks, POS/CMS compliance lookups, and NPI provider summaries, with vector-backed memory, a custom CPT API for accuracy and low latency, and n8n-orchestrated workflows. 

SharkTank1

Table of Contents

Why this matters

What we built

System architecture

Key components

Quick start

Detailed setup

1) Core services: Postgres (pgvector) + n8n

2) Google Vertex AI (Gemini) setup

3) Custom CPT Lookup API (Go)

4) Load n8n workflows

Using the agents

Caching, RAG, and fallbacks

Security & compliance

Project structure

Roadmap

Troubleshooting

Contributors

License

Why this matters

Provider teams depend on the Portal to check prior authorization (PA), place of service (POS) rules, and provider details. Outages, slow responses, or manual multi-step lookups create downstream issues: delayed approvals, coding errors/denials, frustrated users, and compliance risk. Our project targets this friction by automating lookups, unifying data across sources, and delivering a single, reliable answer—even during partial outages. 

SharkTank3

 

SharkTank1

What we built

FrictionFreeAgent is a set of three agentic bots orchestrated by n8n:

Prior Authorization Assistant – Given CPT code(s) and state, determines if PA is required and injects contextual policy signals (e.g., site of service).

Compliance Checker (POS/CMS) – Verifies POS rules and aligns with CMS/UHC policy artifacts.

Provider Search – Summarizes a provider profile from NPI(s), aggregating across internal/external streams and surfacing an accuracy score for name/location alignment. 

SharkTank3

 

SharkTank2

To maximize accuracy and reduce latency, we added a custom CPT Lookup API (Go) and pair it with pgvector memory + Vertex AI (Gemini) for reasoning, summarization, and RAG. 

SharkTank2
