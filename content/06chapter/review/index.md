---
title: "Machine Learning Review"
linktitle: "Machine Learning Review"
weight: 50
---

### Review 

By the end of this exercise, you should observe the following:
	
1. **Anomaly Detection Model Status**  
   - The model progresses through **Collecting**, **Building**, and **Running** stages.  
   - Status shows **Running** before you begin launching attacks.

2. **Traffic Visibility**  
   - Legitimate traffic appears in the **Threat Analytics** dashboard as normal traffic with no detections.  
   - When running **ml-mix**, with attack mix 30%, both legitimate and malicious requests appear in logs.

3. **Attack Detection**  
   - Malicious requests (**SQL Injection**, **Command Injection**, **XSS**) are detected and flagged in the **Attack Logs**.  
   - Log details show **attack type**, **source IP**, and **parameter field** targeted.
4.	**Mitigation Actions**
	- Depending on policy, malicious traffic is either blocked or alerted.
	- Blocked events show the relevant action in log details.


