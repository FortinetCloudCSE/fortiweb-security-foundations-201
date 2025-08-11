---
title: "Task 1: Summary of Bot Mitigation Modules"
linkTitle: "Task 1: Summary of Bot Mitigation Modules"
weight: 10
---


|                            |    |  
|----------------------------| ----
| **Goal**                   | Review Bot Mitigation Modules available in FortiAppsec Cloud.
| **Task**                   | Read the summary provided below. If you need additional information please ask instructor.
| **Verify task completion** | N/A


FortiAppSec Cloud offers multiple bot mitigation techniques, each targeting different aspects of automated threat detection and prevention.

---

### 1. Biometric-Based Detection
Monitors browser interactions such as mouse movements, screen touches, and scroll actions within a set timeframe.  
This helps confirm that requests are coming from a real human rather than an automated process.

---

### 2. Threshold-Based Detection
Lets you define detection rules for suspicious behaviors based on occurrence, time period, severity, and trigger policy.  
Common use cases include:

1. **Crawler Detection** – Identifies excessive crawling behavior that may signal automated indexing or scraping.  
2. **Vulnerability Scanning** – Detects automated scans looking for application or infrastructure weaknesses.  
3. **Slow Attack Detection** – Flags slow-rate requests designed to evade detection or overwhelm servers.  
4. **Content Scraping Detection** – Recognizes automated scraping of web content for unauthorized use.  
5. **Illegal User Scan Detection** – Catches scanning activity aimed at finding and exploiting vulnerabilities.


---

### 3. Bot Deception
Inserts hidden links into HTML response pages.  
Legitimate users never see or click these links, but automated bots often will.  
Requests to these hidden resources are strong indicators of bot activity.

---

### 4. Known Bots
Protects websites, mobile apps, and APIs from both harmful and legitimate bots—covering DoS bots, spam bots, crawlers, and more—without blocking critical automated traffic.  
Includes two predefined rules, plus the ability to create custom ones tailored to your needs.  
Once a Known Bot rule is triggered, its traffic can bypass additional scans.

---

### 5. Machine Learning (ML)-Based Bot Detection
Uses AI-driven detection alongside signature and threshold rules to identify sophisticated bots that might otherwise slip through.  
The model analyzes user behaviors across 13 dimensions, such as request frequency, HTTP version compliance, and resource access patterns.

#### How It Works
FortiAppSec Cloud uses the **Support Vector Machine (SVM)** algorithm to:
- Learn traffic profiles of legitimate clients.
- Compare new client behavior to known patterns.
- Flag anomalies as potential bot traffic.

The process runs in three phases:

**Phase 1 – Sample Collection**  
- Captures behavioral data (samples) during visits.  
- Splits into 75% training data and 25% testing data.

**Phase 2 – Model Building**  
- Analyzes training samples to create behavior profiles.  
- Adjusts SVM parameters to remove outliers and refine accuracy.  
- Selects the best model based on accuracy, cross-validation, and test results.

**Phase 3 – Model Running**  
- Compares incoming traffic to established profiles.  
- Flags significant deviations as anomalies.  
- Triggers actions like alerts or blocking for repeated anomalies.  
- Runs bot confirmation checks to reduce false positives.  
- Updates the model automatically if legitimate traffic patterns change.

---

By combining these techniques—biometric tracking, threshold rules, deception, known bot filtering, and ML-based detection—FortiAppSec Cloud delivers a layered and adaptive defense against automated threats.
