---
title: "Task 2- Setup API Gateway"
menuTitle: "Task 2- Setup API Gateway"
weight: 20
---

1.  From the FortiWeb Cloud Console left pane, select **ADD MODULES**.  Scroll down and turn on **API Gateway** under API Protection.

![api on](api-on.png)

2.  Now API PROTECTION should show up on the left side of the screen. Under API PROTECTION, select **API Gateway**

3.  Add a **Name** and **Email address** Then Click **OK**

![api user](api-user.png)

4.  Next click **Create API Gateway Rule**.  

- for this section, choose a name.  For both "Frontend" and "Backend", enter ```/rest/``` then click **Add URL Prefix**
- turn on API Key Verification
- choose **HTTP Header** for API Key In
- for Header Field Name enter ```apikey```
- for Allow Users, select the user you created in step 3
- leave the Rate limits at default
- select **OK**

![api rule](api-rule.png)

5. Now your user should have an API key click on the eye icon to display the key.  Copy it and put it into a note pad.

![see key](see-key.png)

6. Ensure that the action is set to **Alert & Deny** and then click **Save**

![api save](api-save.png)
