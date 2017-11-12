This program is simple. It grabs the top 5 headlines from either your reddit home page or /all.

In order to make it work you will have to first create a new app on your reddit dev account.

  1. Log in to your reddit account.
  2. Click on the button in the upper right to enter into your account settings.
  3. Click on the applications tab.
  4. Scroll to the bottom and click create another app.
  5. Make sure that you mark it as a script/for personal use.
  6. In the jsontemplate folder there is a file called userinfo.txt.
  7. Fill it out with the relevant information and then save it in the root directory as text.json.

That should configure it to run. Currently there are no binaries supplied. Running it without any arguments supplied
will return the top 5 posts of your homepage. Running it with `all` will return the top 5 posts of /all.
