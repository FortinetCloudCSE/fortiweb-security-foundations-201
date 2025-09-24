---
title: "Task 3: Simulate Bot Traffic"
linkTitle: "Task 3: Simulate Bot Traffic"
weight: 30
---

|                            |    |  
|----------------------------|----
| **Goal**                   | Run traffic to test the FortiAppSec Cloud configuration  
| **Task**                   | Use the provided script to simulate bot traffic  
| **Verify task completion** | If successful, simulated bot traffic will be flagged and blocked in FortiAppSec Cloud  

---

### Use the provided tool to send traffic to your application

We’ve provided a tool called **bots** that will simulate a mix of bot and legitimate traffic.

1. Log into Kali Linux:  ```https://{{Kali IP}}:8443```

2. Open the terminal emulator by clicking the black box icon at the top of the Kali home screen.  
At the prompt, type:  ./bots
    {{< figure src="accessTerminal.png" alt="accessTerminal" >}}


3. You will be prompted with the following menu: ( notice most of the values are default values)

- Load a saved profile: `Y`  
- Target URL: `https://<FortiAppSecStudentID>.fwebtraincse.com`  
- Optimize for OWASP Juice Shop? `Y`  
- Choose number: `6`  
- CSV log file: *(leave blank)*  
- Total run duration: `4m`  
- Concurrency: `20`  
- HTTP timeout per request: `12s`  
- Progress interval: `10`  
- Use one sticky IP in `X-Forwarded-For`: `n`  
- Force HTTP/1.0: `n`  
- Requests per second: `20`  
- Start now: `Y`

{{< figure src="Bot-Attack-1.png" alt="Bot-Attack-1" >}}

Once the tool starts running, you’ll see progress updates approximately every 10 seconds.

{{< figure src="Bot-Attack-2.png" alt="Bot-Attack-2" >}}

The bots tool will send traffic to your application protected by FortiAppSec.
Wait about 4 minutes, then start checking the dashboard and log files in FortiAppSec.


{{< notice >}}
If you lose access to the FortiAppSec Console, open an <strong>Incognito</strong> browser and use the link below to log back in:

<pre><code>https://customersso1.fortinet.com/saml-idp/proxy/demo_sallam_okta/login/</code></pre>
{{< /notice >}}


### check the FortiAppSec Dashboards and Log Files

There are multiple ways to review the logs. We’ll start with the “big picture” view using incidents on the dashboard.

 1.  From the FortiAppSec Cloud Console, select *** Threat Analytics*** from the left-hand menu.
 {{< figure src="Bot-Attack-3-1.png" alt="Threat-Analytics" >}}

 2. In the Top Attack Types pane click on Bot Attacks (scanner) to view incident details.
  Several drill-down options will show source IP, source country, URL attacked, and more.
  You’re encouraged to explore these options. 
 {{< figure src="Bot-Attack-4.png" alt="Incident" >}}
 {{< figure src="Bot-Attack-5.png" alt="Incident-1" >}}
 {{< figure src="Bot-Attack-6.png" alt="Incident-2" >}}

 3. Under Threat Analytics, click on Attack Logs to view individual log entries with detailed information.

 {{< figure src="Bot-Attack-7.png" alt="Access-logs" >}}

we can now look at individual log entries. Logs can be Filtered by Application, URL, Source Country etc. Logs can also be filters by time period. 
please click on any log entry and investigate the detailed information provided. 

{{< figure src="Bot-Attack-8.png" alt="Access-Log-Detail" >}}