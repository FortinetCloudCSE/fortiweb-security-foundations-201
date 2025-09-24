---
title: "Task 2: Configuring Bot Mitigation Modules"
linkTitle: "Task 2: Configuring Bot Mitigation Modules"
weight: 20
---


|                            |    |  
|----------------------------| ----
| **Goal**                   | Setup Bot Mitigation using FortiAppSec Cloud
| **Task**                   | Enable and Configure all the Bot Mitigation Modules on FortiAppSec CLoud
| **Verify task completion** | If successful, Simulated Bot Traffic will be flagged and blocked on FortiAppSec CLoud 

### Enable Bot Mitigation Modules 

{{< notice >}}
If you lose access to the FortiAppSec Console, open an <strong>Incognito</strong> browser and use the link below to log back in:

<pre><code>https://customersso1.fortinet.com/saml-idp/proxy/demo_sallam_okta/login/</code></pre>
{{< /notice >}}

 1.  From the FortiAppSec Cloud Console select your application and in the left pane, select **Waf >** **ADD MODULES**.  Scroll down and turn on **Known Bots, Threshold Based Detection, Biometric Based Detection and Bot Deception** under Bot Mitigation.

   {{< figure src="Bot-1.png" alt="Enable-Bot" >}}

2.  Now the Bot Mitigation tools we need to configure should show up on the left side of the screen. 


#### Configure Known Bots Module 

1. Under **Bot Mitigation** , select **Known Bots**

2. Change the value in the Action box to **Alert and Deny** , **enable** the **known Bad Bots** toggle switch Then Click **Save**

   {{< figure src="Bot-2.png" alt="Known-Bots" >}}

#### Configure Threshold Based Detection Module 

1.  Under **Bot Mitigation** , select **Threshold Based Detection**

2. Change the value in the Action box to **Alert and Deny** , **enable** the **Crawler, Vulnerability Scanning, Slow Attack,Content Scraping,Credential Based Brute Force** toggle switch. Then Click **Save**

{{< figure src="Bot-3.png" alt="Threshold Based" >}}

#### Configure Biometrics Based Detection 

1.  Under **Bot Mitigation** , select **Biometrics Based Detection**

2. click on **Create Rule**, type in **photo** in the URL box, Click **OK** to continue. Then Click **Save**

{{< figure src="Bot-4.png" alt="Biometric" >}}

#### Configure Bot Deception

1.  Under **Bot Mitigation** , select **Bot Deception**

2. click on **Create Rule**, type in **about** in the URL box, Click **OK** to continue. Then Click **Save**

{{< figure src="Bot-5.png" alt="Bot Deception" >}}