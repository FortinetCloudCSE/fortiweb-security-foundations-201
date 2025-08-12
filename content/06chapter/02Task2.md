---
title: "Configure and test the Anomaly Detection module"
linkTitle: "Configure and test the Anomaly Detection module"
weight: 20
---


## Enabling Anomaly Detection in FortiAppSec

In this section, we will enable the **Anomaly Detection** module, which uses machine learning to block zero-day threats and other sophisticated attacks.  
This module builds a behavioral model by analyzing legitimate traffic patterns, allowing it to detect anomalies and unknown attack types.

To train the model, we will use a tool that generates a sufficient number of legitimate requests.  
Please note: the tool may take approximately **30 minutes** to complete its run.

### Step 1: Enable the Anomaly Detection Module

1. From the **FortiAppSec Cloud Console**, select your application.
2. In the left navigation pane, go to **WAF > ADD MODULES**.
3. Scroll down to the **Security Rules** section.
4. Toggle **Anomaly Detection** to **On**.

![Anomaly-on](anomaly-on.png)


In a production environment, both **known attack detection** (signature-based) and **anomaly detection** are used together.  
However, for this demonstration, we will disable signature-based detection.

1. From the **FortiAppSec Cloud Console**, select your application.
2. In the left navigation pane, go to **WAF > Security Rules > Known Attacks**.
3. In the **Signature-Based Detection** pane, disable the following by toggling each corresponding button:  
   - SQL Injection  
   - Cross-Site Scripting  
   - Generic Attacks  
   - Known Exploits  
   - Trojans  
4. Click **Save**.



### Step 2: Run the Tool to Generate Legitimate Traffic

To build the anomaly detection model, you need to generate enough legitimate requests.

1. Open a terminal window from your Kali desktop.

   Run the following command:

   ```./ml-2`` 

2. When prompted, enter the URL you are targeting in the following format:


   ```https://<FortiWebStudentID>.fwebtraincse.com```

   ![start-tool](ML-2.png)

 The tool will begin sending requests to simulate legitimate user traffic.
 You will see progress messages in your terminal indicating how many requests have been sent and the time remaining.
 Note: Progress may take up to 90 seconds to appear — please be patient and avoid starting the tool multiple times. 

![progress](ML-3.png)


⚠️ Note: The process may take up to 30 minutes to complete. Please keep the terminal open during this time.

3. While the tool is running, you can log into the FortiAppSec Cloud Console to observe the model-building progress.


### Step 3: Review the Anomaly detection module on FortiAppSec Cloud

{{< notice >}}
If you lose access to the FortiAppSec Console, open an <strong>Incognito</strong> browser and use the link below to log back in:

<pre><code>https://customersso1.fortinet.com/saml-idp/proxy/demo_sallam_okta/login/</code></pre>
{{< /notice >}}





1. From the **FortiAppSec Cloud Console**, select your application.
2. In the left navigation pane, select **Waf > Security Rules > Anomaly Detection** 
   Click on the TreeView tab and drill down to the search parameter field. You will see the stages: **Collecting, Building, and Running.**

   ![ML_build_1](ML_build_1.png)

   {{< notice warning >}}

Building the model can take up to 30 minutes. Please do not delete once it is built. we will need it for the next exercise 

{{< /notice >}}

Once the tool finishes running, you will see a completion message.


3. When the model reaches the Running stage, you are ready to proceed with launching attacks.

![finished running](ML-4.png)


### Step 4: Launch Attacks 

To test the model, we will run another tool that launches **SQL Injection, Command Injection**, and **XSS attacks** along with legitimate requests..

1. Open a terminal window from your Kali desktop.

   Run the following command:

   ```./ml-mix``` 

2. When prompted, enter the URL:


   ```https://<FortiWebStudentID>.fwebtraincse.com```

Accept the default values for the remaining options:
 
  - Duration: 5m 
  - Target 30
  - Workers: 20
  - Attack mix percentage: 30 
  - Use /rest/products/search? q= ... : Y ( this is the parameter field we trained in the pervious section)
  - Skip TLS verification: **Y**
  - Per- request timeoue : 10s
  - Verbose sample logging: n 

While the tool is running, log back into the FortiAppSec Cloud Console and review the logs to confirm that attacks are being detected and mitigated.
