FrictionFreeAgent 🚀

Agentic AI to Reduce Operational Friction in Healthcare

📌 Overview

FrictionFreeAgent is an agentic AI proof-of-concept designed during a summer internship project to improve operational efficiency in healthcare provider portals. When UnitedHealthcare’s Prior Authorization and Notification (PAAN) application faces downtime, errors, or delays, providers experience friction that leads to increased call volumes, claim denials, compliance risks, and diminished user trust
.

Our solution leverages agentic AI, RAG pipelines, and workflow automation to proactively minimize this friction by providing reliable, accurate, and real-time support.

⚡ Problem Statement

Healthcare provider portals play a critical role in:

Cost control – preventing unnecessary or expensive treatments.

Regulatory compliance – ensuring procedures align with guidelines.

Patient safety – validating clinical appropriateness.

Operational efficiency – reducing downstream costs
.

However, service outages, manual provider lookups, and mismatched data create inefficiencies:

Delayed prior authorizations.

Incorrect billing/coding → claim denials.

Frustrated providers and overwhelmed support teams.

Compliance risks from outdated information
.

💡 Solution: FrictionFreeAgent

We designed FrictionFreeAgent, a modular multi-agent AI system, to:

Automate and validate provider lookups across multiple sources.

Provide unified, consistent data from structured (APIs, databases) and unstructured (PDFs, documents) inputs.

Eliminate errors and reduce manual dependency.

Enhance user experience with conversational AI
.

Key Capabilities

Prior Authorization Assistant

Determines if prior authorization (PA) is required for a given CPT/procedure code.

Provides contextual details by state and service location
.

Compliance Checker (POS Lookup)

Cross-checks Place of Service (POS) with CMS and UHC guidelines.

Ensures compliance across internal and public sources.

Provider Search Agent

Aggregates provider details from multiple APIs (internal and public).

Scores accuracy by comparing data across streams
.

🛠️ Tech Stack & Tools

LLMs: Google Gemini models (via Vertex AI) for conversational AI, summarization, and structured outputs
.

Workflow Automation: n8n
 for orchestrating multi-agent pipelines and API integrations
.

Vector Storage: Postgres PGVector for embedding and retrieving structured healthcare data (CPT, CMS, POS, NPI)
.

Custom APIs: FastAPI-based CPT lookup service to ensure accuracy and reduce latency
.

Data Sources:

Internal Optum/UHC APIs (NPI, CPT, PA rules).

Public APIs (NIH NPI Registry, CMS).

Locally stored data for redundancy and training.

🧩 Architecture
flowchart TD
    A[Provider Input (CPT, NPI, POS)] --> B[Workflow Orchestration - n8n]
    B --> C1[Custom FastAPI CPT Lookup API]
    B --> C2[Postgres PGVector Store]
    B --> C3[Internal + Public APIs]
    C1 --> D[Gemini Model - Vertex AI]
    C2 --> D
    C3 --> D
    D --> E[AI-Powered Chat UI]
    E --> F[Unified Response with Accuracy Score]


Chat-based interface replaces manual form lookups.

AI agents fetch, validate, and summarize provider/authorization info.

Accuracy scoring ensures transparency for providers
.

📊 Business Impact

Reduced downtime impact – system serves as a reliable backup during outages.

Lower claim denials – improved accuracy in PA and provider info.

Operational efficiency – fewer manual lookups and reduced support team load.

Improved provider experience – proactive notifications, empathetic messaging, and seamless chatbot interactions
.

🔮 Roadmap & Next Steps

API Access Expansion – integrate dynamic CMS/POS APIs instead of static PDFs
.

Multi-Agent Autonomy – implement advanced agents:

Outage Notifier Agent

Task Resubmission Agent

Backpressure Agent

Model Training & Feedback Loop – use locally stored provider portal interactions to fine-tune Gemini models for accuracy.

Production Deployment – CI/CD pipelines, UAT testing, monitoring dashboards, and security reviews.

📽️ Demo

👉 [Add link here if you have a hosted demo, Loom recording, or slides]

📜 License

This project was developed as part of a summer internship at Optum (UnitedHealth Group).
All rights reserved © 2025 Optum, Inc.
