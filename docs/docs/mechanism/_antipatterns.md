# Antipatterns 
Several behaviours exist on the exchange that are detrimental to the success of the openlab exchange and the overall LabDAO community. We list some of these recurring patterns below and how we limit their occurrence. 

## wash trading 
LabDAO might distribute governance tokens to users of the openlab exchange. To reduce the amount of unproductive trading, we take the following actions:
* the same input and output address can only be used when the value of the transaction is 0 - some labs might want to use this approach when they have a service in house to keep track of all their data. 
* the fees on the exchange will be high enough to make it unattractive to "mine" the LAB token by washtrading. The DAO might choose to distribute tokens proportional to the value of unique requests.
* wallets that are involved in wash trading might lose their status as a trusted provider.
