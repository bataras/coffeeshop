Coffee Shop Simulator in GO
=======
I drank a lot of coffee while coding this coffee shop simulator. I should code a microbrewery simulator next...

Everyone knows what a coffee shop does, more or less. Customers line up to place orders.
Someone takes those orders with the customer's name and lines them up on the counter or otherwise hands them off.
Baristas run around making the coffee and calling out the customer's name and order when it's ready.
Customers wait to hear their name, grab their coffee and either walk out or sit and drink for a while.

Some coffee shops are big with multiple workers and can fulfill several orders at the same time.
Some are single person operations where the barista does everything more or less sequentially.

All coffee shops have machines that can do discreet things in the background and alert or be recognized in some
way that they're done. For example, grinders grind for a while and then stop. Brewing machines run hot water through
the ground coffee and then stop. Meanwhile, the baristas drive things in the foreground, interacting with the
customers, operating the machines or doing other tasks like maintenance and cleaning.

We can decompose this into a simplified model as follows...

How to Run
--------
**Build and run**
```
go build cmd/coffeeshop/main.go
./main
```

**Get help**
```
./main -h

Usage of config:
  -baristas int
    	number of baristas. -1 means use config file (default -1)
  -conf string
    	path to a .yaml config file (default "coffeeshop.yaml")
  -customers int
    	number of customers. -1 means use config file (default -1)
```

**Configure**

Please see config.yaml. It's pretty self-explanatory with comments.

See the results of orders...
```
./main >out
grep "complete. took" out
time="03:46:00.7474" level=error msg="--- Order 4: complete. took 95.794µs: no grinder pool for: italian"
time="03:46:01.14892" level=error msg="--- Order 6: complete. took 135.366µs: no grinder pool for: italian"
time="03:46:03.91426" level=info msg="--- Order 5: complete. took 2.965789169s"
time="03:46:03.91432" level=error msg="--- Order 10: complete. took 190.395µs: no grinder pool for: italian"
time="03:46:05.46731" level=info msg="--- Order 7: complete. took 4.117534529s"
time="03:46:05.97251" level=info msg="--- Order 8: complete. took 4.421452261s"
time="03:46:06.84365" level=info msg="--- Order 12: complete. took 2.527784843s"
time="03:46:07.37389" level=info msg="--- Order 13: complete. took 2.857144459s"
time="03:46:08.22001" level=info msg="--- Order 9: complete. took 4.506646346s"
time="03:46:08.26912" level=info msg="--- Order 11: complete. took 4.154146254s"
```

See what a single barista did...
```
./main >out
grep "Barista 1" out
time="03:46:00.74724" level=info msg="Barista 1: handle cash register with 2 orders in the pipe"
time="03:46:00.94832" level=info msg="Barista 1: took order Order#: 5 Beans: french Ounces: 9 Strength: NormalStrength"
time="03:46:00.94861" level=info msg="Barista 1: grind start Order#: 5 Beans: french Ounces: 9 Strength: NormalStrength"
time="03:46:02.84999" level=info msg="Barista 1: grind start Order#: 8 Beans: columbian Ounces: 9 Strength: MediumStrength"
time="03:46:04.85171" level=info msg="Barista 1: grind start Order#: 13 Beans: ethiopian Ounces: 14 Strength: NormalStrength"
time="03:46:05.97245" level=info msg="Barista 1: brewer done. give coffee to customer Order#: 8 Beans: columbian Ounces: 9 Strength: MediumStrength"
time="03:46:05.97255" level=info msg="Barista 1: grinder refill bean: ethiopian hopper 48 48"
time="03:46:06.53475" level=info msg="Barista 1: grind start Order#: 11 Beans: columbian Ounces: 8 Strength: MediumStrength"
time="03:46:08.26903" level=info msg="Barista 1: brewer done. give coffee to customer Order#: 11 Beans: columbian Ounces: 8 Strength: MediumStrength"
```

There can be much more metrics gathering; like how busy each barista, grinder and brewer are

Primary Entities and Resources
=========

Sketches. Stuff changed after these. But it shows early thoughts.

[Early sketch 1](sketch-1.jpg)

[Early sketch 2](sketch-2.jpg)

Beans
--------
- has a type (Columbian, Ethiopian, French, Italian, etc)
- beans are created out of thin air
- there is an endless supply of them
- the model is adaptable to later add a "Roaster" from which the coffee shop orders more beans

Extraction Profiles
--------
- matches an ExtractionStrength to a grams-per-ounce conversion
- could be updated to to do things like grinder settings: drip, espresso, etc
- a real coffee shop will occasionally adjust extraction profiles during the day based on humidity, etc.
  - but for POC all we do is convert strength to grams/ounce

Grinder
--------
- has a hopper filled with a single type of beans
- hopper can run out and be refilled in H(G) time, where G is how many grams are needed
- can only be refilled with the same type of beans
- can grind X grams of beans in grinder-specific G(X) time
- only grinds for one order at a time
- can only be used by one barista at a time
- never has to be cleaned or breaks down
- baristas cannot do anything else while waiting for the grinder to grind

Brewer
--------
- can pass X ounces of hot water through any quantity of beans in B time
- has an endless supply of hot water
- only brews for one order at a time
- can only be used by one barista at a time
- never has to be cleaned or breaks down
- because brewing can be slow, baristas can do other things while waiting
  - barista A can start brewing, then take a customer's order while barista B fulfills A's order

Customer
--------
- waits for a barista to take her order
- placing the order takes a fixed amount of time P
- requests a coffee by size, bean type and strength
- size is 8, 12 or 16oz
- strength is light, medium, strong
  - regulated internally by the coffee shop as a certain amount of beans to be ground
- waits for the order to complete
- does not have to disambiguate her order completion from other customer orders
- leaves the shop and doesn't sit down

Barista
--------
- waits for a customer who wants to order
- takes the order
- fills grinders
- grinds the beans
- brews the beans
- completes the order
- can do other tasks while waiting for a brewer to complete
  - including completing another barista's order when it's done brewing
- can start and manage multiple brewers in parallel as load/resources demand
- the model is adaptable to later add behaviors like taking breaks, cleaning, leveling up on skills, etc

Cash Register
--------
- a construct where a customer and a barista meet and take a little time to create an order and put in in the queue

Order
------
- specifies size (small, medium, large... 8, 12 or 16oz)
- specifies strength (light, medium, strong)
- specifies bean type
- contends for grinders by bean type and availability
- contends for brewers by availability

Coffee Shop Configuration
-------
- has 1 or more baristas
- has 1 or more grinders
- has 1 or more brewers
- has 1 or more types of beans
- coffee strengths: grams of beans needed for light, medium, strong
- coffee sizes: small, medium, large... 8, 12 or 16oz
- P time for customer to place order
- grinder-specific refill and grind rates
- brewer-specific brew rates

Operational Ideals
-------
- customers shouldn't generally wait to be able to place an order
  - but as a rule of thumb, there should never be more that twice the number of orders in the pipeline than there are baristas to fulfill them
- baristas should never be idle unless there is nothing to do
- orders should be fulfilled as a FIFO
  - but may not always be due to grinder/bean type contention
  - or parallel baristas with faster grinders/brewers
- orders should be fulfilled as fast as possible
  - example: it's better to refill grinder hoppers without adding to the customer wait time. ie, when both the grinder and a barista are idle
- a barista should not dequeue an order in isolation and then wait for the corresponding grinder to be available
  - there may be a subsequent order that can be processed instead of waiting for a grinder for the prior order
- baristas fundamentally start processes and then wait for them to complete (eg a brewer)
  - baristas time-slice across customers and resources in the coffee shop
    - for example A can start brewing, then take a waiting customer X's order while barista B fulfills A's order after grinding beans for customer Y

Operational Metrics
---------
Ran out of time to add these. But it would be fairly straight forward given the code design.
Was fantasizing about docker composing a Prometheus and Grafana system to show graphs.

- Customer satisfaction APDEX score
- avg barista idle time
- grinder idle times
- brewer idle times
- grinder and brewer parallelism
- ad nauseam...

Order Processing Pipeline
===========

Resource Pools
-------
- one availability pool for each bean type specific hopper/grinder
- one availability pool for the brewers

Processing Queues
-------
- one customer waiting to order mechanism (cash register)
- one order priority queue by order time
- one queue for grinder hopper refill requests (pointer to the grinder)
- one availability queue for grinders based on the bean type they hold
- one priority queue for each paired order/grinder. priority by order time
- one queue for brewer machines that are done brewing

Order/Grinder Pairing
----------
- pairs available grinders with oldest matching order
  - newer orders can potentially process before older ones depending on grinder/bean availability
- forwards matched pairs to the order/grinder priority queue based on order age

Barista Handling Logic
----------
- baristas loop and listen to...
  - the brewer availability queue
  - the order/grinder priority queue
  - the grinder hopper refill request queue
  - the customer waiting to order queue
- if a brewer is done
  - complete the order back to the original customer
  - return the brewer 
- if # of orders in the pipeline is <= 2 x the number of baristas and a customer is waiting
  - take customer order
  - reject if bean type not available in the coffee shop
- if a grinder hopper needs to be refilled
  - refill the hopper unless another barista got to it first (including one who is using it for an order)
  - this provides an opportunity for hopper refills to not cause a longer wait for a customer. but not guaranteed
- if an order/grinder pair is available
  - fill the grinder hopper if needed
  - grind the beans
  - if the hopper is X% empty, send a refill request to the grinder hopper refill request queue
  - return the grinder to the availability pool
  - wait for a brewer to be available (a signal)
  - start the brewer
    - have the brewer add itself to the brewer done queue when it's done
- loop forever


Tech Choices
==========
General approach...
- Use google, youtube, awesome-go (https://awesome-go.com/)
- Find libraries on github
- Look for...
  - popularity (stars, forks, views etc)
  - recency: when was the last commit
  - stability: is the author still playing with it, fixing/testing stuff
  - authority: are tech influencers talking about or using it
- Read the code and examples
  - is it clean
  - obvious errors (I've seen this)
  - is it trying to do too much, overkill
- Future proof-y
  - is the work being done as POC
  - does the library's interface appear to be generalized and abstracted enough to reasonably extensible
  - does it feel too tightly coupled for the usage model
  - on the other hand, sometimes it's ok to knowingly adopt something with which you can grow longer term
  - balancing tech debt
