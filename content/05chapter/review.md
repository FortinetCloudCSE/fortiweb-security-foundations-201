---
title: "API Protection Review"
linktitle: "API Protection Review"
weight: 50
---

### Review 

In this module, we onboarded our application into FortiWeb Cloud using the GUI.  We also used the built in firewall rules in our GCP Virtual Private Cloud Network to lock down access to the origin server such that only traffic from FortiWeb Cloud will be accepted.


### API Protection Quiz

1. What is the final step for onboarding a Web Application in FortiWeb Cloud?

    {{% expand title="Click here for answer" %}}
Change the DNS Record.  While we did not perform this step for the purposes of this lab.  In a production environment, the final step to onboarding your application is to change either the CNAME or A record for your application such that all traffic is directed towards FortiWeb Cloud.
    {{% /expand %}}

2. You must use TLS on port 443 to communicate from FortiWeb Cloud to your origin server. (True or False)

    {{% expand title="Click here for answer" %}}
False: While it is highly recommended to use TLS for the connection from FortiWeb Cloud to the origin server, as we saw in the lab, the server protocol and port are configurable.
    {{% /expand %}}

3. Web attacks are difficult to perpetrate and you need to be an inveterate hacker to attempt it (True or False)

    {{% expand title="Click here for answer" %}}
**FALSE** - The attack in this lab is very simple, but very effective.  This should highlight the need to protect web applications with a purpose built Web Application Firewall (WAF)
    {{% /expand %}}