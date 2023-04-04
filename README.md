Coffee Shop Simulator in GO

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

Beans
--------
- has a type (Columbian, Ethiopian, French, Italian, etc)
- beans are created out of thin air
- there is an endless supply of them
- the model is adaptable to later add a "Roaster" from which the coffee shop orders more beans

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

Ad-hoc Grinder
--------
- has no hopper
- must be loaded with X grams of beans L(X) time prior to grinding
- can grind X grams of any type of bean in grinder-specific G(X) time
- only grinds for one order at a time
- can only be used by one barista at a time
- never has to be cleaned or breaks down
- because grinding is quick, baristas cannot do anything else while waiting for the grinder to grind

Brewer
--------
- can pass X ounces of hot water through any quantity of beans in B time
- has an endless supply of hot water
- only brews for one order at a time
- can only be used by one barista at a time
- never has to be cleaned or breaks down
- because brewing can be slow, baristas can do other things while waiting

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
- instantly conjures the required amount of beans into existence as needed
- fills grinders
- grinds the beans
- brews the beans
- completes the order
- can do other tasks while waiting for a brewer to complete
  - including completing another barista's order when it's done brewing
- can start and manage multiple brewers in parallel as load/resources demand
- the model is adaptable to later add behaviors like taking breaks, cleaning, leveling up on skills, etc

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
- if it has > 1 bean types B, it must have either B hopper-grinders or at least 1 ad-hoc grinder
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
  - this implies an order/grinder pairing stage ensuring that a barista doesn't work on an order until there is a grinder for it

Operational Metrics
---------
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
- one availability pool for each ad-hoc grinder
- one availability pool for each brewer

Processing Queues
-------
- one customer waiting to order queue
- one order queue
- one queue for grinder hopper refill requests (pointer to the grinder)
  - lower priority than orders, or higher ?
- ephemeral resource-affinity queues for matching orders with grinders
- one priority queue for each paired order/grinder. priority by order time
- one queue for brewer machines are done brewing

Order/Grinder Pairing Stage
----------
- pairs available grinders with oldest matching order
  - using ephemeral resource-affinity queues that are created on demand
  - newer orders can potentially process before older ones depending on grinder availability
- forwards matched pairs to the order/grinder priority queue based on order age

Barista Handling Stage
----------
- baristas loop and listen to...
  - if currently handling an order/grinder job from the previous iteration
    - the brewer availability queue
  - otherwise
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
- if an order/grinder pair from the Pairing Stage is available
  - fill the grinder hopper if needed
  - grind the beans
  - if the hopper is X% empty, send a refill request to the grinder hopper refill request queue
  - return the grinder to the availability pool
  - wait for a brewer to be available
  - start the brewer
    - have the brewer add itself to the brewer done queue when it's done. use a closure probably
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
