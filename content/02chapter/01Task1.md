---
title: "Task 1: Onboard Application"
linkTitle: "Task 1: Onboard Application"
weight: 10
---

|                            |    |  
|----------------------------| ----
| **Goal**                   | Start protecting Juice Shop Application with FortiAppSec Cloud
| **Task**                   | Onboard Application in FortiAppSec GUI
| **Verify task completion** | Your Application will show up in the Application list.

### Add Application

1. Click on the **WAF** icon in the left menu bar, open the **Applications** view, and then click **+ ADD APPLICATION**.

    ![Add-Application1](Add-application.png). \
        \
        \
        and then click, **+ ADD APPLICATION** 

    ![Add-Application2](Add-application-2.png). 
           
          
2. **_Tab 1: "WEBSITE"_** 

   - In **Web Application Name** enter your FortiAppSec Cloud StudentID number which you used to login to FortiAppSec Cloud (found at the top right corner of the FortiAppSec Cloud Screen).   

    {{% notice info %}}For example, if your FortiAppSec Cloud User is **CSEAccount669@fortinetcloud.onmicrosoft.com**, your Student ID would be: **669**{{% /notice %}}

   - For **Domain Name** use ```<studentId>.fwebtraincse.com``` and then select next
    
    ![conf-app1](conf-app1.png)

2. **Tab 2: Network**,

   - Clear the **"HTTP"** as we want to force users to interact with FortiAppSec using only HTTPS.
   - For **IP Address or FQDN** enter the **JuiceShop Public IP** (which is the Ubuntu VM Public IP from your Terraform Output)
   - For **Port** enter "3000"
   - Select **HTTP** for Server Protocol.  This is Juice Shop and it is NOT secure
   - Click on **Test Origin Server**  You should see a green box pop up that says "Test successfully"
   - Choose **Next**
    
    ![Conf-app2](conf-app2.png)

3. **Tab 3: CDN** 

    **_No Changes_**.  You will notice the Selected WAF Region shows the Platform "AWS" and the Region. In your lab it may show a different platform and region  
    
    {{% notice info %}}FortiAppSec Cloud automatically chooses the platform and region based on the IP Address of the application.  There is no user intervention required.{{% /notice %}}
    
    - Select **Next**
    
    ![conf-app3](conf-app3.png)

4. **Tab 4: "SETTING"**

   - **DO NOT** enable Block Mode

   - Select **Save**
   
    ![conf-app4](conf-app4.png)

5. **Tab 5: "CHANGE DNS"**

   We are presented with very important information regarding DNS settings which need to be changed in order to direct traffic to FortiAppSec Cloud.  In this lab, we will not be doing this, as sometimes it can take a while for the DNS settings to propagate.  

   {{% notice warning %}} 
   Take Note of the IPv4 addresses and CNAME for use in a later step.  **Before you close!**
   {{% /notice %}}

   - Select **Close**
   ![conf-app5](conf-app5.png)

6. You should now see your Application listed in FortiAppSec Cloud.  Note that the DNS Status is set to **Update Pending** This is expected, and we will ignore it.
    ![conf-app6](conf-app6.png)

   {{% notice note %}} If you need to recover the application IPs or CNAME later, you can click on the app's DNS status **Update Pending** to show DNS status & retrieve the IPs

   ![](app-ips.png)
   {{% /notice %}}

7. Update Google DNS

    {{< dns_record >}}
