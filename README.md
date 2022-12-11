# Galt - Salt but better... I hope

The whole idea behind Galt is to make Salt but better, faster and hopefully less prone to obscurities such as the order of things to run, "include" things before the include statement etc. Aka create a better version of something great, but fix a few things.

## How is Galt supposed to work

Galt is supposed to be scalable to thusands of servers but not use much resources like Salt that is basically edge triggered aka the Salt master dictates how Salt minions do things where as Galt will be level triggered aka the desired state is defined and the Galt minion will do what it can to make it so. This means that the Galt minion will not be told "when to apply something" it will watch for desired state change in the Galt master and then do what it can to do so, using templates etc. etc. just like Salt.

# Word of warning

Even though I've been working with computers and work as a sysadmin and been at it for +2decades then do remember that this project will in no way guarantee any safety of your systems at this point - not take any responsibility if you delete anything on your minions - it is provided as is.

