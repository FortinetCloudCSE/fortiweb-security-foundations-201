---
title: "Task 3 - Explore FortiWeb Options"
linkTitle: "Task 3 - Explore FortiWeb Options"
weight: 30
---


FortiWeb Cloud Options

In the previous task, we simply turned on Block Mode in FortiWeb Cloud.  This enabled the default, minimum security configuration.  Take a moment now to click through some of the menu options on the left to see what Features are enabled by default.  We will also look at how to enable new features.

1. Navigate to **Security Rules** on the left menu and click on **Known Attacks** to see what features are turned on.  The first category is Signature Based Detection.  Click the **Search Signature** button on the right and search for the injection Keyword.  

![Search-Sig](search-sig.png)

2. On the left menu, click through the available menus for **Access Rules, Bot Mitigation and DDOS Prevention**

3. **Vulnerability Scan** is an additional paid service that can be added to FortiWeb Cloud, which will scan your protected Applications for OWASP Top 10 vulnerabilities.

{{% notice info %}}
More information can be found in the docs at:
https://docs.fortinet.com/document/fortiweb-cloud/23.3.0/user-guide/898181/vulnerability-scan
{{% /notice %}}

4. Next Click on **+ Add Modules**.  This is where we can activate additional security featuresfeatures.  These features are all covered under the FortiWeb Cloud WAF-as-a-Service License, which is billed based on the number of websites protected and the average Mbps throughput in aggregate for all protected sites.

{{% notice info %}}
FortiWeb Cloud Datasheet:
https://www.fortinet.com/content/dam/fortinet/assets/data-sheets/fortiweb-cloud.pdf
{{% /notice %}}
