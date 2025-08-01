---
title: "Task 3: CSRF attack "
linkTitle: "Task 3: CSRF attack "
weight: 20
---

|                            |    |  
|----------------------------| ----
| **Goal**                   | Use Burpsuite to Perform CSRF attack
| **Task**                   | Burpsuite will modify a user password, using Cross Site Request Forgery, then block it with FortiWeb Cloud
| **Verify task completion** | The CSRF password change should go through initially, but be blocked after enabling CSRF protection on FortiWeb Cloud

A Cross-Site Request Forgery (CSRF) attack is a type of security exploit where an attacker tricks a user into performing actions on a web application without their consent. This can happen when a malicious website, email, or other online resource causes the user's web browser to perform an unwanted action on a different site where the user is authenticated.

1. Let's generate a CSRF attack with Burpsuite. 

2. Repeat Step 1-5 from Task 1 to open Burpsuite. if Burpsuite is already running in the background just click to go back to at by clicking on the top left corner of Kali linux.

3. On the proxy tab, Click on **Open Browser**

    ![csrf1](csrf1.png)

4. Type the FQDN allocated: ```https://<studentId>.fwebtraincse.com``` into the browser.

    ![csrf2](csrf2.png)

5. Once the Juiceshop app loads, click on Account > Login.
    {{% notice note %}} If you don't see **Account** in top right bar, you may have to expand the browser window {{% /notice %}}
    ![csrf3](csrf3.png)

6. Create a user login by clicking on **Not Yet a customer?** at the bottom. 

    ![csrf4](csrf4.png)

7. Make sure to use the same email and credentials as below just so we won't forget. 

   - email: ```test@example.com```
   - password: ```test1234$```
   - Repeat Password: ```test1234$```
   - Security Question: Select **Your eldest sibling's middle name** from dropdown. 
   - Answer: ```botman```

   - Click on **register**

    ![csrf5](csrf5.png)

8. Login using the credentials above. 

    ![csrf6](csrf6.png)

    ![csrf7](csrf7.png)

9. Once logged in clik on Account > Privacy and Security > Change Password. 

   - Current password: ```test1234$```
   - New Password: ```password1234$```
   - Repeat New Password: ```password1234$```

   - Click **Change**

    ![csrf8](csrf8.png)

10. Once changed we can see **your password was successfully changed** dialog. 

    ![csrf9](csrf9.png)

11. Go back to **Burpsuite > Proxy > HTTP History** and Scroll down to the end to see the last HTTP call made which is the **/rest/user/change-password**. Right-click on the change-password GET call and select **send to repeater**. 

    ![csrf10](csrf10.png)
    
    ![csrf11](csrf11.png)

12. Click on the **Repeater** tab to see the change password request. The Raw request shows the current password and new password we updated. 

    ![csrf12](csrf12.png)

13. Execute a Cross Site Request Forgery password change attack!
    - Remove the current password field from the request
    - Update the request to reflect only new and repeat password using: ```hello1234$```
    - Your request should look like below:
    - Click **Send** after the request is updated. 
        ![csrf13](csrf13.png)

16. Response is a 200 OK meaning that call is successful. 

    ![csrf14](csrf14.png)

17. Verify by going back to juiceshop, account login. Logout if already logged in. 

    ![csrf15](csrf15.png)

18. Account > login 

    - email: ```test@example.com```
    - password: ```hello1234$```
    - Click **Log In**

    ![csrf16](csrf16.png)

    As we can see with successfully login using the new credentials, our CSRF attack was successful!

19. Now login to **FortiAppSec Cloud**
    - Be sure to click on your **allocated application.**

20. Scroll down to **Waf** **->** **Add modules** at the bottom. Add **CSRF protection** under Client Security Module and click **OK**

    ![Add-Module](Add-Module.png)
    
    ![add CSRF](<Add CSRF.png>)

21. In the Application View > WAF > Client Security > Click on CSRF Protection.
    - On both **Page List Table** AND **URL List Table**, Add the URL ```/rest/user/change-password```
    - Update the Action to **Alert and Deny** and click Save. the Module takes ~3 minutes to be in effect. 

    ![CSRF-config](CSRF-config.png)
    
    ![CSRF-config-1](CSRF-config-1.png)


22. Once done, repeat the attack again with Password of your choice, and you should see a block message. 

    ![csrf21](csrf21.png)

23. On Fortiappsec cloud, **Threat Analytics > Attack Logs >** There is a CSRF attack log.

    ![csrf22](csrf22.png)
