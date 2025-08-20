# FrictionFreeAgent — Agentic AI to Reduce Operational Friction in Healthcare

> Public showcase of a summer internship project that uses **agentic AI**, **RAG**, and **workflow automation** to keep provider operations moving when portals are slow or unavailable.  
> Built with **Google Gemini (Vertex AI)**, **n8n**, **Postgres + pgvector**, and a **FastAPI** microservice.

---

### Table of Contents
- [Why](#why)
- [What It Does](#what-it-does)
- [System Architecture](#system-architecture)
- [Key Components](#key-components)
- [Quick Start](#quick-start)
- [Configuration](#configuration)
- [RAG & Accuracy Scoring](#rag--accuracy-scoring)
- [Security & Compliance Readiness](#security--compliance-readiness)
- [Roadmap](#roadmap)
- [Screenshots / Demos](#screenshots--demos)
- [Team](#team)
- [License](#license)

---

### Why

Provider portals (e.g., **PAAN — Prior Authorization and Notification**) are mission-critical for cost control, compliance, patient safety, and operational efficiency. When outages, latency, or data mismatches occur, clinicians and staff face delays, denials, and frustration:contentReference[oaicite:0]{index=0}:contentReference[oaicite:1]{index=1}.  

**FrictionFreeAgent** addresses this by **automating lookups**, **unifying data across sources**, and **surfacing consistent answers via a chat UI**, even when underlying systems are degraded:contentReference[oaicite:2]{index=2}:contentReference[oaicite:3]{index=3}.

---

### What It Does

**FrictionFreeAgent** is a modular, multi-agent system:

- **Prior Authorization Assistant**  
  Determines whether **PA** is required for a CPT/procedure code and adds state/POS context:contentReference[oaicite:4]{index=4}:contentReference[oaicite:5]{index=5}.

- **Compliance Checker (POS)**  
  Cross-checks **Place of Service** against internal UHC policy and public **CMS** sources to flag discrepancies:contentReference[oaicite:6]{index=6}.

- **Provider Search Agent (NPI)**  
  Aggregates provider details from **internal** and **public** APIs, returns a single structured summary, and assigns an **accuracy score** by comparing sources:contentReference[oaicite:7]{index=7}:contentReference[oaicite:8]{index=8}.

**Why this design?**  
- Automates and validates lookups to reduce manual error and latency:contentReference[oaicite:9]{index=9}  
- Unifies structured & unstructured data (APIs, PDFs, databases) for consistency:contentReference[oaicite:10]{index=10}  
- Conversational UX replaces brittle form flows and speeds task completion:contentReference[oaicite:11]{index=11}  

---

### System Architecture

```mermaid
flowchart TD
    U[User / Provider] -->|CPT, NPI, POS prompts| UI[Chat UI]
    UI --> N8N[n8n Orchestrator]

    N8N --> FA[FastAPI CPT Lookup]
    N8N --> VDB[(Postgres + pgvector)]
    N8N --> APIs[Public + Internal APIs]
    subgraph External Sources
      NIH[NIH NPI Registry]
      CMS[CMS / POS Data]
    end
    APIs --> NIH & CMS

    FA --> LLM[Gemini on Vertex AI]
    VDB --> LLM
    APIs --> LLM

    LLM --> RESP[Unified Structured Response + Accuracy Score]
    RESP --> UI
