---
title: "Positive security model and machine learning"
linkTitle: "Positive security model and machine learning"
weight: 10
---


The **positive security model**, also known as **whitelisting**, operates on the principle of allowing only known good behavior while blocking everything else. In this approach, the system defines a set of explicitly permitted actions—such as valid URLs, parameters, and methods—and denies all other actions by default. This contrasts with the **negative security model**, which focuses on blocking known bad behavior.

Positive security models are effective at preventing both known and unknown attacks because they explicitly define what is allowed, reducing the attack surface and providing strong protection against unauthorized access or malicious activity. However, maintaining and updating whitelists can be challenging, especially in dynamic environments where applications and usage patterns change frequently.

To address these challenges, **FortiAppSec** uses **machine learning for anomaly detection**. Its anomaly detection model monitors URLs, parameters, and HTTP methods of HTTP and/or HTTPS sessions targeting your web applications, building mathematical models to detect abnormal traffic behavior.

## Two Layers of Machine Learning in FortiAppSec

**FortiAppSec** employs two layers of machine learning to detect malicious attacks:

1. **First Layer: Hidden Markov Model (HMM)**  
   - Monitors application access.
   - Collects data to build mathematical models for every parameter and HTTP method.
   - Evaluates each request against the learned model to detect anomalies.

2. **Second Layer: Threat Model Verification**  
   - If the first layer flags a request as anomalous, the second layer determines if it is a true attack or a benign anomaly.
   - Uses pre-built, pre-trained **threat models** for categories like SQL Injection, Cross-site Scripting (XSS), etc.
   - Models are trained using thousands of attack samples and continuously updated via the **FortiAppSec Security Service**.
   - The **FortiGuard** team analyzes new threats and retrains relevant models, which are then pushed to all FortiAppSec installations similarly to signature updates.

---

## How FortiAppSec Builds Its Anomaly Detection Model

**FortiAppSec** constructs its machine learning model by evaluating domain-specific parameters based on extensive samples of legitimate requests.

### Sampling Criteria

A request is treated as a sample if **all** of the following are true:

- The response code is `200` or `302`.
- The response `Content-Type` is `text` or `html`.
- The request includes parameters in the URL or body.

### Pattern Generalization

When a valid sample is collected, **FortiAppSec** generalizes it into a pattern.  
For example:
- `"abcd_123@abc.com"` and `"abcdefgecdf_12345678@efg.com"`  
  → generalized to → `"A_N@A.A"`

The model is built on **patterns**, not raw values.

### Model Lifecycle

- **Initial Model:**  
  Created after collecting **400 samples**.  
  Actively used to detect anomalies while more data is collected.

- **Model Promotion:**  
  Once **1200 samples** are collected, the system evaluates the stability of the patterns:
  
  - **Stable Patterns:**  
    If few new patterns are seen, the model is promoted to a **standard model**.
  
  - **Unstable Patterns:**  
    If many new patterns are still emerging, sample collection continues until stability is achieved.

- **Standard Model:**  
  More accurate and reliable.  
  Continuously updated as application behavior evolves (e.g., new URLs or parameter changes).  
  Outdated patterns are discarded, and new patterns are introduced to keep the model current.

---

By leveraging this multi-layered machine learning approach, **FortiAppSec** provides robust protection against both known and emerging web application threats.
