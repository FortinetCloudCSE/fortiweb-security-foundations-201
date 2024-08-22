---
title: "Task 2- Enable Block Mode"
menuTitle: "Task 2- Enable Block Mode"
weight: 20
---


1. Enable Block Mode on FortiWeb Cloud

On the Applications page enable block mode by clicking on the Block Mode button

![En-Block](en-block.png)

2. Repeat the same step to perform SQLi attack in the browser.

```sh

For example: https://669.fwebtraincse.com/?name=' OR 'x'='x

```

{{% notice info %}}
You will see that FortiWeb now blocks the SQLi attack.
{{% /notice %}}

![Blocked](blocked.png)


3. Now let's navigate to our application page in FortiWeb Cloud, by clicking on the Application Name.  This should take you to the Application **Dashboard**.  You should see a Threat listed in the **OWASP Top 10 Threats box called A03:2021-Injection.  Click on it.

![App-dash](app-dash.png)

4. Navigate through some of the tabs.

![INJ-Det](inj-det.png)

5. On the **Threats** tab, click on the Threat.  In this case **Known Attacks**.  This will take you to a list showing dates when this type of attack was encountered.  If you click on the Arrow next to the date, more information about that incident can be seen.  Spend some time clicking around on the Clickable links in this output.  There is a lot of information available from here, including a link to the OWASP Top 10 site describing this attack as well as HTTP header information and matched patterns.

![KA-Det](ka-det.png)