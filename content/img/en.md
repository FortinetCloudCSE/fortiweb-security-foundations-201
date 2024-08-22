# FortiWeb
FortiWeb is Fortinet's Web Application and API security platform, enabling enterprise customers to protect web applications no matter where they are deployed.  FortiWeb defends web applications and APIs against OWASP Top-10 threats, DDOS attacks, and malicious bot attacks. Advanced ML-powered features improve security and reduce administrative overhead. Capabilities include anomaly detection, API discovery/protection, bot mitigation and advanced threat analytics to identify the most critical threats across all protected applications.

FortiWeb Cloud WAF as a Service (WAFaaS) reduces administrative overhead by offering the full suite of FortiWeb security features without the need to manage VMs and networking.  Customers simply change their DNS records so that all appication traffic is proxied through FortiWeb Cloud.

## Overview
In this lab, you will have an opportunity to configure FortiWeb Cloud to protect a Juice Shop server, emulating a very vulnerable e-commerce website.  Students will onboard their application and then run a simple SQL injection attack and an SQL map attempt.  As the format for these training session are purposefully short, in order to maintain student interest and fit into the work day.  The attacks used here are very high level. We will not be covering the ML feature set in this lab, as the training period requirements would be time prohibitive.


### Objectives
In this lab:

- When you start this lab, two VMs will be created in two different GCP projects.  One of these is Kali Linux, which will be used to attack the other, which is Juice Shop.  

## Setup and requirements
### Before you click the Start Lab button

<ql-warningbox>
Read these instructions.
</ql-warningbox>

Labs are timed and you cannot pause them. The timer, which starts when you click **Start Lab**, shows how long Google Cloud resources will be made available to you.

You should have recieved an email from **fortinetsecdevops@gmail.com** with a subject of **ForiWeb_Cloud_Training**.  This email has the proper credentials required to log into FortiWeb Cloud.  Links to the appropriate login and logout pages will be provided below.

This hands-on lab lets you do the lab activities yourself in a real cloud environment, not in a simulation or demo environment. It does so by giving you new, temporary credentials that you use to sign in and access Google Cloud for the duration of the lab.

To complete this lab, you need:

* Access to a standard internet browser (Chrome browser recommended).  
    >*Note: Use an Incognito or private browser window to run this lab. This prevents any conflicts between your personal account and the Student account, which may cause extra charges incurred to your personal account.*

* Time to complete the lab---remember, once you start, you cannot pause a lab.  
> *Note: If you already have your own personal Google Cloud account or project, do not use it for this lab to avoid extra charges to your account.*

### How to start your lab and sign in to the Google Cloud Console

1. Click the **Start Lab** button. If you need to pay for the lab, a pop-up opens for you to select your payment method. On the left is the **Lab Details** panel with the following:
    * Time remaining
    * Your temporary credentials that you must use for this lab
    * Your temporary project ID
    * Links to additional student resources

2. Open Google Cloud console in new **incognito** browser window by right clicking the **Google Cloud Console** link in **Student Resources**.
    ***Tip:*** Arrange the tabs in separate windows, side-by-side.
    >*Note: If you see the Choose an account dialog, click Use Another Account.*

3. Copy the **GCP Username** and **Password** from the **Lab Details** panel and paste it into the Sign in dialog. Click **Next**.
    > Important: You must use the credentials from the left panel. Do not use your Google Cloud Skills Boost credentials.

    >*Note: Using your own Google Cloud account for this lab may incur extra charges.*

4. Click through the subsequent pages:
    * Accept the terms and conditions.
    * Do not add recovery options or two-factor authentication (because this is a temporary account).
    * Do not sign up for free trials.

## Lab Environment

Below is a diagram of the Lab environment.

![lab1](./img/diagram.png)

### Task 1: Check Availability of Juice Shop

Use the public IP of Juice Shop to log in: http://{{{protected_project.startup_script.Juice-Shop-IP | Juice Shop IP}}}:3000

![Juiceshop Home Page](./img/juice-home.png)

### Task 2: Login to Kali

1.  Use the Kali public IP to log in: https://{{{kali_project.startup_script.Kali-IP | Kali IP}}}/vnc.html

Accept certificate errors and proceed.  When prompted, click **Connect**.  This will take you to the home screen of Kali

![Kali Home Page](./img/kali-home.png)

2.  In order to copy/paste into Kali, we will need to click on the tab at the left hand side of the screen.

![cp-tab](./img/cp-tab-kali.png)

3.  This will open the tab revealing a couple of options.  Select the clipboard icon and paste your text into the box.  Once the text is in the box, you can right click on the desktop and select Paste Selection to paste in the text.  When done, you can click on the arrow to hide the clipboard.

![paste-kali](./img/paste-kali.png)

4.  We are going to make a small change to Kali in order to prepare for later steps. Please open a terminal by clicking on the Black Icon at the bottom of the screen.  Then enter the following:

```sh
bash
echo nameserver 8.8.8.8 >> /etc/resolv.conf
wget https://dl.pstmn.io/download/latest/linux_64 -O /tmp/linux_64 && tar xvzf /tmp/linux_64 -C /tmp/ && sudo mv /tmp/Postman /opt/ && sudo ln -s /opt/Postman/app/Postman /usr/local/bin/Postman
```

### Task 3: Log Into FortiWeb Cloud

1.  Using an **Incognito** browser, navigate to the below URL:

```sh

https://customersso1.fortinet.com/saml-idp/proxy/demo_sallam_okta/login 

```

2.  Input the username from the email you recieved from **fortinetsecdevops@gmail.com** and click **Next**

![FWeb login](./img/fweb-login.png)

3.  Input the password from the email you recieved from **fortinetsecdevops@gmail.com** and click **Sign in**

![FWeb pass](./img/fweb-pass.png)

For the next step, choose **Yes**.  You do want to stay logged in.

<ql-infobox> Sometimes if you wait too long to input your password, you will get SAML login portal error "Error: SAML response with InResponseTo is too late for previous request"  If this happens just click the small blue "Login" link. </ql-infobox>

4. This will take you to the FortiCloud Premium Dashboard. At the top of the screen select **Services** > **FortiWeb Cloud**

<ql-infobox>When you log in, you will see that you are unauthorized to view the FortiCloud Premium Dashboard.  This is expected, as this user has not been given this permission.</ql-infobox>

![unauth](./img/unauthorized.png)

![Choose fweb](./img/choose-fweb.png)

## Onboard Juice Shop Application

### Task 1: Add Application

1. At the top of the screen, clci, on **+ ADD APPLICATION** 

![Add-App](./img/add-app.png)

2. For Step 1 "WEBSITE" 

- for **Web Application Name** enter the number of the username found in the email you recieved from **fortinetsecdevops@gmail.com**.   

<ql-infobox>For example, if the Username is CSEAccount669@fortinetcloud.onmicrosoft.com the number would be 669</ql-infobox>

- For **Domain Name** use number.fwebtraincse.com and then select next

![App-1](./img/app-1.png)

2. For Step 2,

- **unselect "HTTP"** as we want to force users to interact with FortiWeb using only HTTPS.
- For **IP Address or FQDN** enter the JuiceShop Public IP {{{protected_project.startup_script.Juice-Shop-IP | Juice Shop IP}}}
- For **Port** enter "3000"
- Select HTTP for Server Protocol.  This is Juice Shop and it is NOT secure
- Click on **Test Origin Server**  You should see a green box pop up that says "Test successfully"
- Choose **Next**

![App-2](./img/app-2.png)

3. For Step 3 "CDN" we will not change anything.  You will notice the Selected WAF Region shows the Platform "Google Cloud Platform" and the Region.  

<ql-infobox>FortiWeb Cloud automatically chooses the platform and region based on the IP Address of the application.  There is no user intervention required.</ql-infobox>

- Select **Next**

![App-3](./img/app-3.png)

4. In Step 4 "SETTING" we will **NOT** enable Block Mode

- Select **Save**

![App-4](./img/app-4.png)

5. In Step 5 "CHANGE DNS" We are presented with very important information regarding DNS settings which need to be changed in order to direct traffic to FortiWeb Cloud.  In this lab, we will not be doing this, as sometimes it can take a while for the DNS settings to propagate.  

<ql-warningbox> 
Take Note of the IPv4 addresses and CNAME for use in a later step.  **Before you close!**
</ql-warningbox>

- Select **Close**

![App-5](./img/app-5.png)

6. You should now see your Application listed in FortiWeb Cloud.  Note that the DNS Status is set to **Update Pending** This is expected and we will ignore it.

![App-on](./img/app-on.png)

<ql-warningbox>This is a **Shared Environment** !!!  Please ensure that you are only making changes to **Your Application**.  After Applications are onbaorded into FortiWeb Cloud, Administrators have full RBAC capabilities, but we will not be activating that during this 90 minute lab.</ql-warningbox>

## Secure Google Infrastructure

This lab was originally bootstrapped with and ingress firewall rule which allows all ports and protocols from all sources (0.0.0.0/0).  This is not a best practice.  It is recommended, to only allow necessary ports and/or sources.  Now that we have onboarded our application, we want to ensure that the only device that can communicate with our application is FortiWeb Cloud.

<ql-infobox>
For the below steps, ensure that you are in the Protected Project in the GCP Console.  When you log into the console for the first time, you are directed to the Protected Project by default.  This can be changed by clicking on the drop down at the top left of the console screen, between the Google Logo and the Search bar.
</ql-infobox>

Protected Project ID = {{{protected_project.project_id | Protected Project}}}

### Task 1: Modify GCP Network Firewall Rules

1. First, let's grab the IPs which FortiWeb Cloud will use to communicate with your application. From the FortiWeb Cloud Applications page, select **Allow IP List** from the top of the page.  This will open a dialog showing Management and Scrubbing Center Addresses.  Clcik on **Copy to Clipboard**.  Paste these IPs into a text document and then click **Return**

![WAF-IP](./img/waf-ip.png)

2. In the Google Console, select the "Hamburger" menu at the top left of the screen and navigate to **VPC network > Firewall**

![Hamburger](./img/hamburger.png)

3. Click on the Ingress rule named **terraform-fweb-qwiklab-untrust-vpc-(random)-ingress**

![G-POL](./img/g-pol.png)

4. Click **EDIT** at the top of the page.

5. Scroll down to **Source IPv4 ranges**.  Delete **0.0.0.0/0** and enter the IP's copied from FortiWeb Cloud in step 1.  Then Click Save

![ED-FW](./img/ed-fw.png)

6. To ensure that it worked, use your browser to try and navigate to Juice Shop: http://{{{protected_project.startup_script.Juice-Shop-IP | Juice Shop IP}}}:3000

You should **NOT** be able to access Juice Shop Directly.

## Simple SQL Injection

### Task 1: Perform a simple SQL injection attack

According to the Open Worldwide Application Security Project (OWASP):

<ql-infobox>
A SQL injection attack consists of insertion or “injection” of a SQL query via the input data from the client to the application. A successful SQL injection exploit can read sensitive data from the database, modify database data (Insert/Update/Delete), execute administration operations on the database (such as shutdown the DBMS), recover the content of a given file present on the DBMS file system and in some cases issue commands to the operating system. SQL injection attacks are a type of injection attack, in which SQL commands are injected into data-plane input in order to affect the execution of predefined SQL commands.
</ql-infobox>

You can find more information at "https://owasp.org/www-community/attacks/SQL_Injection"

For this task, we will just use a simple Browser.

1. Log into Kali linux: https://{{{kali_project.startup_script.Kali-IP | Kali IP}}}/vnc.html

2. BSince we did not modify the DNS record we will enter a host entry on Kali.  Open the terminal emulator by clicking on the black box at the bottom of the Kali Hom screen. At the prompt, type:

```sh
bash
sudo nano /etc/hosts
```

3. When the host file opens you will need to add in the host entry in the format "ip address number.fwebtraincse.com" at the bottom of the file.  For this you can enter one or both of the IP Addresses you noted earliear during "Step 5" of the onboarding earlier.   Once this is complete, type **ctrl+o** followed by **enter** and then **ctrl+x**

![Hosts](./img/hosts.png)

4. Now let’s Navigate to the browser (located at thebottom of Kali home page) and type the URL. into the navigation bar https://number.fwebtraincse.com.  Accept warnings and proceed to the application

<ql-infobox>If we had modified the DNS record at the begining of this lab, FortiWeb would have pulled a valid SSL certificate from Let's Encrypt.</ql-infobox>

![stud-home](./img/stud-home.png)

5. Let’s perform a very simple SQLi attack. To perform a SQLi attack append ?name=' OR 'x'='x to your URL.  Be sure that you use **YOUR NUMBER**.  Below is just an example.

```sh

For example: https://669.fwebtraincse.com/?name=' OR 'x'='x

```

<ql-infobox>
The attack will go through and you will see the Juice Shop Home page
</ql-infobox>

### Task 2:  Protect WebServer from Attack

1. Enable Block Mode on FortiWeb Cloud

On the Applications page enable block mode by clicking on the Block Mode button

![En-Block](./img/en-block.png)

2. Repeat the same step to perform SQLi attack in the browser.

```sh

For example: https://669.fwebtraincse.com/?name=' OR 'x'='x

```

<ql-infobox>
You will see that FortiWeb now blocks the SQLi attack.
</ql-infobox>

![Blocked](./img/blocked.png)


3. Now let's navigate to our application page in FortiWeb Cloud, by clicking on the Application Name.  This should take you to the Application **Dashboard**.  You should see a Threat listed in the **OWASP Top 10 Threats box called A03:2021-Injection.  Click on it.

![App-dash](./img/app-dash.png)

4. Navigate through some of the tabs.

![INJ-Det](./img/inj-det.png)

5. On the **Threats** tab, click on the Threat.  In this case **Known Attacks**.  This will take you to a list showing dates when this type of attack was encountered.  If you click on the Arrow next to the date, more information about that incident can be seen.  Spend some time clicking around on the Clickable links in this output.  There is a lot of information available from here, including a link to the OWASP Top 10 site describing this attack as well as HTTP header information and matched patterns.

![KA-Det](./img/ka-det.png)

### Task 3:  FortiWeb Cloud Options

In the previous task, we simply turned on Block Mode in FortiWeb Cloud.  This enabled the default, minimum security configuration.  Take a moment now to click through some of the menu options on the left to see what Features are enabled by default.  We will also look at how to enable new features.

1. Navigate to **Security Rules** on the left menu and click on **Known Attacks** to see what features are turned on.  The first category is Signature Based Detection.  Click the **Search Signature** button on the right and search for the injection Keyword.  

![Search-Sig](./img/search-sig.png)

2. On the left menu, click through the available menus for **Access Rules, Bot Mitigation and DDOS Prevention**

3. **Vulnerability Scan** is an additional paid service that can be added to FortiWeb Cloud, which will scan your protected Applications for OWASP Top 10 vulnerabilities.

<ql-infobox>
More information can be found in the docs at:
https://docs.fortinet.com/document/fortiweb-cloud/23.3.0/user-guide/898181/vulnerability-scan
</ql-infobox>

4. Next Click on **+ Add Modules**.  This is where we can activate additional security featuresfeatures.  These features are all covered under the FortiWeb Cloud WAF-as-a-Service License, which is billed based on the number of websites protected and the average Mbps throughput in aggregate for all protected sites.

<ql-infobox>
FortiWeb Cloud Datasheet:
https://www.fortinet.com/content/dam/fortinet/assets/data-sheets/fortiweb-cloud.pdf
</ql-infobox>

## Dig Deeper

Now that we have done a simple SQL injection attack, let's take a deeper dive into one of the tools that an actual hacker (or Red Team) might actually use to attack an application.

### Task 1: Use Burp Suite to find a vulnerability

Burp Suite gives us a quick and easy way to query targeted sites.

1. At the bottom of the Kali home page, click on the terminal icon (black box).  Once open, input:

```sh

burpsuite

```

2. Burp Suite will pop up. Accept all of the warnings and EULAs.  Leave Temporary Project selected and click **Next**

![Burp_Suite1](./img/bs-temp.png)

3. Leave "Use Burp defaults" selected and click **Start Burp**.

![Burp_Suite2](./img/bs-start.png)

4. Accept the warning that Burp Suite is out of date and then select settings at the top right of the screen.

![Burp_Suite3](./img/bs-set.png)

5. In the settings menu, select **Burp's browser**.  Under **Browser running** check the box for "Allow Burp's browser without a sandbox"

![BS-sand](./img/bs-sand.png)

>*Note: once the button is clicked, just close the settings menu.  There is no need to save.*

6. Click on the **Proxy** tab at the top of the Burp Suite screen.  This will bring you to the Intercept screen.  Click on **Open Browser**.  Ignore Error and proceed to next step.

![Burp_Suite5](./img/bs-proxy.png)

7. In the browser URL bar, input https://number.fwebtraincse.com and hit enter.  This will bring you to the juice shop home page.

8. Minimize the browser and click on the **HTTP History** tab under Proxy.  Scroll down the list until you find a URL labeled **"/rest/products/search?q=**.  Select this line and right click.  Then click on **Send to Repeater**.  This will allow us to manipulate the requests in order to do a little nefarious recon.

![BS-URL](./img/bs-url.png)

9. At the top of Burp Suite, Click on the **Repeater** Tab.  You will see the request we just sent.  Now click on the **Send** Button.  This will populate the Response area.

![Burp_Suite8](./img/bs-repeater1.png)

10. Now we are going to modify our query a bit.  Click on the First line in the Raw request and add **'--** to our get request after.  The GET should now look like **/rest/products/search?q='--**.  Click **Send**.  We will now see an error in the Response section.  This error tells us that the database is SQLITE and uncovers a vulnerability.

![Burp_Suite9](./img/bs-repeater2.png)

<ql-infobox>It's worth mentioning that the standard signature based Web Protection Profile did not catch this attempt. If Machine Learning were enabled, this would not have gone through.  Instead it would have been identfied as an anomaly and then passed to the threat engine where it would have been identified as an SQL Injection attempt.  We are not using ML in this lab, as the number of samples required to train the Model would be time prohibitive</ql-infobox>

### Task 2: Use SQLMAP to exploit vulnerability

Now that we know what the Database type is, we can use sqlmap to see if we can get some "Juicy" information (pun intended).  You could just run SQLMAP initially to find the vulnerability, but It would take much longer without an idea of what you were looking for.

1. Open a new terminal on Kali, and take a look at the SQLmap help page.  I also think it's helpful to use bash shell here, as we will want to be able to use the up arrow in order to scroll though old commands

```sh

bash
sqlmap -h

```

2. Now we will attempt to discover what typ SQL injection vulnerabilities exist.  Since we know that the database runs on **sqlite** we can shorten the scan time by giving sqlmap that information.  Input the first line below at the terminal, substituting your URL.

```sh

sqlmap -u "https://number.fwebtraincse.com/rest/products/search?q="  --dbms=SQLite --technique=B --level 3 --batch

```

<ql-infobox>This attempt will fail, due to the default protections offered by FortiWeb.  It is still recommended to use ML in production in order to prevent reconnasiance from previous step</ql-infobox>

![Map-Blocked](./img/map-blocked.png)

3. Disable Block Mode on your application in FortiWeb Cloud

![Dis-Block](./img/dis-block.png)

4. Re-run the sqlmap attempt.  You will see that some vulnerabilities were found.

![Map-Allow](./img/map-allow.png)


## Task 3: CSRF attack 


A Cross-Site Request Forgery (CSRF) attack is a type of security exploit where an attacker tricks a user into performing actions on a web application without their consent. This can happen when a malicious website, email, or other online resource causes the user's web browser to perform an unwanted action on a different site where the user is authenticated.

1. Lets generate a CSRF attack with Burpsuite. 

2. Repeate Step 1-5 from Task 1 to open Burpsuite. if Burpsuite is already running in the background just click to go back to at by clicking on the top left corner of Kali linux.

3. On the proxy tab, Click on **Open Browser**

![csrf1](./img/csrf1.png)

4. Type the FQDN allocated: **https://yournumber.fwebtraincse.com** into the browser. for example: **https://670.fwetraincse.com**

![csrf2](./img/csrf2.png)

5. Cnce the Juiceshop app loads, click on Account > Login.

![csrf3](./img/csrf3.png)

6. Create a user login by clicking on **Not Yet a customer?** at the bottom. 

![csrf4](./img/csrf4.png)

7. Make sure to use the same email and credentials as below just so we wont forget. 

- email: ```test@example.com```
- password: ```test1234$```
- Repeat Password: ```test1234$```
- Security Question: Select **Your eldest siblins middle name** from dropdown. 
- Answer: ```botman```

- Click on **register**

![csrf5](./img/csrf5.png)

8. let's login using the credentials above. 

![csrf6](./img/csrf6.png)
![csrf7](./img/csrf7.png)

9. Once logged in clik on Account > Privacy and Security > Change Password. 

- Current password: ```thelloest1234$```
- New Password: ```password123$```
- Repeat New Password: ```password1234$```

- Click **Change**

![csrf8](./img/csrf8.png)

10. Once changed we can see **your password was successfully changed** dialog. 

![csrf9](./img/csrf9.png)

11. Go back to Burpsuite > Proxy > HTTP History and Scroll down to the end to see the latest HTTP call made which is the /rest/user/change-password. Right click on the change-password GET call and select **send to repeater**. 

![csrf10](./img/csrf10.png)

![csrf11](./img/csrf11.png)

12. Click on the repeater tab to see the change password request. The Raw request shows the current password and new password we updated. 

![csrf12](./img/csrf12.png)

13. Remove the current password to only have new and repeat password and retry sending the request to Juiceshop. 

14. Update the request to reflect only new and repeat password as shown below with password to be ```hello1234$```

![csrf13](./img/csrf13.png)

15. Click **Send** after the request is updated. 

16. Response is a 200 OK meaning that call is successful. 

![csrf14](./img/csrf14.png)

17. Verify by going back to juiceshop, account login. Logout if already logged in. 

![csrf15](./img/csrf15.png)

18. Account > login 

- email: ```test@example.com```
- password: ```hello1234$```
- Click **Log In**

![csrf16](./img/csrf16.png)

it Works and its a successful CSRF attack!.

19. Now login to Fortiweb cloud. Make sure to click on your allocated application.

20. Scroll down to **Add modules** at the bottom. Add CSRF protection under Client Security Module and click **OK**

![csrf17](./img/csrf17.png)

![csrf18](./img/csrf18.png)

21. In the Application View > Client Security > Click on CSRF Protection. In the Page list and URL list we will add the URL **/rest/user/change-password** , update the Action to **Alert and Deny** and click Save. the Module takes ~3 minutes to be in effect. 

![csrf19](./img/csrf19.png)

![csrf20](./img/csrf20.png)


22. Once done, repeat the attack again with Password of your choice and you should see a block message. 

![csrf21](./img/csrf21.png)

23. On FortiWeb cloud, Threat Analytics > Attack Logs > There is a CSRF attack log.

![csrf22](./img/csrf22.png)


## FortiWeb Cloud API Gateway

Now that we have run a few attacks, let's explore how to protect our API by adding API Gateway functionality to FortiWeb Cloud.

### Task 1: Call API with Postman

1.  Open postman by opening a new terminal (not bash) and type ```Postman``` at the prompt.  This should start the postman application.

![postman](./img/postman.png)

- When postman comes up, select "Or continue with lightweight API client

![postmanlite](./img/p-light.png)

2.  Now, let's make an api call to search for Apple Juice.  Use the below url, but replace the url with your student url.

```sh
https://<number.fwebtraincse.com>/rest/products/search?q=Apple
```

3.  This first call will fail, due to a certificate error.  In the response section, you will need to scroll down and select "Disable SSL Verification".

![postman ssl disable](./img/p-dis.png)

4. Now the Call should go through an you should see a status 200 and returned data.

![postman success](./img/p-success.png)

### Task 2: Setup API Gateway

1.  From the FortiWeb Cloud Console left pane, select **ADD MODULES**.  Scroll down and turn on **API Gateway** under API Protection.

![api on](./img/api-on.png)

2.  Now API PROTECTION should show up on the left side of the screen. Under API PROTECTION, select **API Gateway**

3.  Add a **Name** and **Email address** Then Click **OK**

![api user](./img/api-user.png)

4.  Next click **Create API Gateway Rule**.  

- for this section, choose a name.  For both "Frontend" and "Backend", enter ```/rest/``` then click **Add URL Prefix**
- turn on API Key Verification
- choose **HTTP Header** for API Key In
- for Header Field Name enter ```apikey```
- for Allow Users, select the user you created in step 3
- leave the Rate limits at default
- select **OK**

![api rule](./img/api-rule.png)

5. Now your user should have an API key click on the eye icon to display the key.  Copy it and put it into a note pad.

![see key](./img/see-key.png)

6. Ensure that the action is set to **Alert & Deny** and then click **Save**

![api save](./img/api-save.png)

### Task 3: Test API Gateway

1.  In Postman, click **Send** again to re-test your api call.  It should return status 403 and return a long error page ending with "Please contact the administrator..."

![no key](./img/no-key.png)

2. Now, let's add a key

- select **Headers** under the URL bar.
- enter ```apikey``` for Key
- enter the previously copied key for Value
- click the empty box next to apikey to send this header
- click **Send**

You should see code 200 and returned data.

![yes key](./img/yes-key.png)

## Task 3: Open API Validation/Schema protection

In this task, lets run through the open API/Swagger based schema protection with Fortiweb cloud. Swagger, now known as the OpenAPI Specification (OAS), is a framework for API development that allows developers to design, build, document, and consume RESTful web services.

example of Swagger: https://petstore.swagger.io/

FortiWeb can validate incoming requests against your OpenAPI schema to ensure they conform to the defined structure and data types. This helps prevent injection attacks and other malicious activities.

1. Download the juiceshop schema file to your local machine by clicking on URL below.

```sh
https://juiceshopswagger.blob.core.windows.net/juiceshopswagger/swagger.yaml?sp=r&st=2024-08-06T16:05:20Z&se=2024-11-09T01:05:20Z&spr=https&sv=2022-11-02&sr=b&sig=F8TWuKSH430782%2FgJBWLhCQuEDK2101CChRkXx4XdU0%3D
```

2. From the FortiWeb Cloud Console left pane, select ADD MODULES. Scroll down and turn on  under API Protection to add OPEN API VALIDATION

![apischema1](./img/api-schema1.png)

3. In the API protection module, click on Open API validation > Create OpenAPI validation file. 

![apischema2](./img/api-schema2.png)

4. Click on "choose file" to upload the file downloaded in Step 1, Click OK. 

![apischema3](./img/api-schema3.png)

5. Dont forget to Save at the bottom. 

![apischema4](./img/api-schema4.png)

6. Now lets open POSTMAN, by opening a new terminal (not bash) and type ```Postman``` at the prompt.  This should start the postman application.

- When postman comes up, select "Or continue with lightweight API client"

![postmanlite](./img/p-light.png)

7. we will send a POST request to the URL we have documented in Schema. Change "**GET**" to "**POST**", for URL use: **https://yournumber.fwebtraincse.com/b2b/v2/orders**

for Request body, Click on Body > Raw > JSON and paste the following:

```sh
{
  "cid": "testing",
  "orderLines": [
    {
      "productId": "testing",
      "quantity": 500,
      "customerReference": 1
    }
  ],
  "orderLinesData": "[{\"productId\": 12,\"quantity\": 10000,\"customerReference\": [\"PO0000001.2\", \"SM20180105|042\"],\"couponCode\": \"pes[Bh.u*t\"},{\"productId\": 13,\"quantity\": 2000,\"customerReference\": \"PO0000003.4\"}]"
}
```
![apischema6](./img/api-schema6.png)

- Note: The schema for Product ID is changed from Integer to String. the Fortiweb cloud Juiceshop schema we uploaded have this value defined as Integer. 

![apischema10](./img/api-schema10.png)

- Click on "**SEND**"

8. We will see "403 internal server error" with a Fortiweb cloud block message in HTML.

![apischema7](./img/api-schema7.png)

9. on Fortiweb Cloud > Attack log > we can see a log generated for this block request to show the reason for block is Open API schema Violation. 

![apischema8](./img/api-schema8.png)

## Delete your Application

1. You are almost done!  Please take a moment to delete only **Your Application** using the trashcan Icon on the right side of the applicaion listing.

![del-app](./img/del-app.png)

2. Please use the below link to log out of FortiCloud

```sh
https://customersso1.fortinet.com/saml-idp/proxy/demo_sallam_okta/saml 
```
Be sure to click the small blue **Logout** button at the bottom of the text.

### Congratulations

Congratulations, you have successfully completed this lab!  Your environment will automatically delete itself at the end of the alloted lab time.
