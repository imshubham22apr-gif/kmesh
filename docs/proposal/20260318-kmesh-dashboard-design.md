---
title: Kmesh Dashboard Design Proposal
authors:
- "@imshubham22apr-gif"
reviewers:
- "@LiZhenCheng9527"
- TBD
approvers:
- "@LiZhenCheng9527"
- TBD

creation-date: 2026-03-18

---

## Kmesh Dashboard Design Proposal

### Summary

This proposal outlines the design and implementation of a unified Kmesh Dashboard. The goal is to lower the usage threshold for Kmesh by providing a graphical interface for common operations such as Waypoint lifecycle management, traffic visualization, and configuration of L7 features like circuit breaking and rate limiting. The dashboard will consist of a React/Next.js frontend and a lightweight Go backend that interacts with Kubernetes via `client-go`.

### Motivation

Kmesh offers advanced traffic management capabilities (L3/L4 via eBPF and L7 via dual-mode/waypoint). However, configuring these features currently requires deep knowledge of Kmesh CRDs and manual YAML manipulation, which can be error-prone and overwhelming for new users. A dashboard will provide a more intuitive, guided experience, making Kmesh more accessible and easier to operate.

#### Goals

- Provide a user-friendly interface for Kmesh management.
- Simplify Waypoint proxy installation and lifecycle management.
- Integrate or embed service topology maps (e.g., Kiali-like visualization).
- Provide low-code wizards for Circuit Breaker and Rate Limit configurations.
- Display Kmesh health status and essential metrics.
- Ensure the dashboard is lightweight and integrates seamlessly with existing Kmesh components.

#### Non-Goals

- Replacing existing Kmesh control plane logic (the dashboard is a consumer/manager of the control plane).
- Providing exhaustive Kubernetes management (focused specifically on Kmesh-related resources).
- Real-time packet inspection (focused on L7 policy and traffic flow).

### Proposal

#### User Stories

##### Story 1: Simplified Waypoint Creation
As a platform engineer, I want to create a Waypoint proxy for a specific namespace or service via a simple form, so I don't have to manually craft Gateway/HTTPRoute resources.

##### Story 2: Visualizing Traffic Health
As an SRE, I want to see a visual map of my services and their traffic status (errors, latency) directly within the Kmesh dashboard to quickly identify issues.

##### Story 3: Configuring Traffic Policies
As a developer, I want to enable a circuit breaker on my service by selecting options in a wizard rather than writing complex Envoy-based configurations.

#### Notes/Constraints/Caveats

- The backend must be lightweight to minimize resource overhead in the cluster.
- The dashboard should support both "Dual Mode" and "Kernel-Native" mode visualizations.
- Authentication and Authorization (RBAC) must be integrated to ensure secure access.

#### Risks and Mitigations

- **Complexity**: Building a full-stack dashboard from scratch is a large undertaking. Mitigation: Phase the development, starting with core Waypoint management and health status.
- **Security**: Direct interaction with K8s API via the dashboard backend. Mitigation: Use refined RBAC roles for the dashboard service account.
- **Maintenance**: Maintaining a separate frontend stack. Mitigation: Use standard, well-documented frameworks (React/Next.js) and maintain clean component abstractions.

### Design Details

#### Architecture Overview

1.  **Frontend**: Based on **React** and **Next.js**.
    - Responsive UI for desktop and tablet.
    - Component-based architecture for reusability (using a modern UI library like Headless UI or Radix).
    - Client-side state management for fluid interactions.
2.  **Backend**: Written in **Go** using **Gin** or **Chi**.
    - Lightweight HTTP server.
    - Uses `client-go` to interact with Kubernetes API server.
    - Interacts with Kmesh specific resources (Waypoints, L7 configurations).
    - Optional: In-memory cache for frequently accessed discovery data.
3.  **Integration**:
    - **Kiali Integration**: Explore embedding Kiali views or consuming Kiali APIs for topology maps.
    - **Monitored Metrics**: Integration with Prometheus/Grafana to display throughput and error rates.

#### Feature Breakdown

- **Waypoint Management**:
  - Dashboard card showing active Waypoints.
  - Form to create Waypoints for Pods, Services, or Namespaces.
- **Traffic Policy Wizards**:
  - Abstractions for Circuit Breaker and Rate Limiting.
  - Validation of inputs before applying CRDs.
- **Mode Toggle**:
  - Visual indicator and toggle for Dual Mode vs Kernel-Native mode for specific workloads.
- **Cluster Health**:
  - Status of Kmesh daemonsets, control plane, and sidecars/gateways.

#### Test Plan

- **Unit Testing**: Go backend logic (API handlers, K8s interaction helpers) and React components/hooks.
- **Integration Testing**: Testing the backend against a local `kind` cluster with Kmesh installed.
- **E2E Testing**: Using Playwright or Cypress to verify critical user journeys (e.g., creating a Waypoint via the UI).

### Alternatives

1.  **CLI Only**: Continue using `kmeshctl` and manual YAML. *Rejected* as it doesn't solve the "usage threshold" problem for less experienced users.
2.  **Generic K8s UI (e.g., Lens extension)**: *Rejected* as it lacks the specialized domain knowledge required for Kmesh's unique dual-mode architecture and traffic flow.
3.  **Direct Envoy UI**: *Rejected* as Kmesh is moving towards a more native/ebpf-integrated model where Envoy is only part of the story (Dual Mode).
