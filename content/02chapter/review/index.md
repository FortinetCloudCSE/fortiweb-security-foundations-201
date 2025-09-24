---
title: "Application Protection Review"
linktitle: "Application Protection Review"
weight: 50
---

### Review 

In this module, we onboarded our application into FortiAppSec Cloud using the GUI.  We also used the built-in Azure Network Security Group (NSG) in our Azure Environment to lock down access to the origin server such that only traffic from FortiAppSec Cloud will be accepted.


### Application Protection Quiz

1. What is the final step for onboarding a Web Application in FortiAppSec Cloud?
    {{% expand title="Click here for answer" %}}
Change the DNS Record.  We performed this in Chapter 2 Task 1 Step 8.  Likewise, in a production environment, the final step to onboarding your application is to change either the CNAME or A record for your application such that all traffic is directed towards FortiAppSec Cloud.
    {{% /expand %}}

2. You must use TLS on port 443 to communicate from FortiAppSec Cloud to your origin server. (True or False)

    {{% expand title="Click here for answer" %}}
**False**: While it is highly recommended to use TLS for the connection from FortiAppSec Cloud to the origin server, as we saw in the lab, the server protocol and port are configurable.
    {{% /expand %}}

3. Why are we no longer able to browse directly to Juice Shop App

    {{% expand title="Click here for answer" %}}
We modified the Azure Network Security group applied to the Juice Shop VM, only allowing FortiAppSec Cloud source IP addresses
    {{% /expand %}}