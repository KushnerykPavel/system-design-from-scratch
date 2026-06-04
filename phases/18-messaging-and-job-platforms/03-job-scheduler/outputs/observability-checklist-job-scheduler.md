# Observability Checklist — Job Scheduler

## Core Metrics

- dispatch lag by priority class
- overdue jobs and missed runs
- retry rate and retry amplification ratio
- shard-scan duration and lease churn
- dead-letter volume by tenant and job family

## Logs

- job creation, pause, resume, and deletion
- manual replay or catch-up actions
- policy changes for retry or deadline behavior

## Review Questions

- Can operators distinguish due jobs from late jobs?
- Can they see whether backlog comes from retries or fresh work?
- Can they identify one tenant causing shard saturation?
