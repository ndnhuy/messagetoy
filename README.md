This is a toy project to learn about messaging system. I intent to build a library of messaging system which hides implementation details of messaging technology such as kafka, rabbitmq, ...

This library is intended to build a simple pub/sub application without knowledge about any underlying techonology. We can even use go channel or database as a event queue as long as it conforms to the interfaces. 

And in order to do that, I need to create an abstraction of messaging components such as Message Channel, Message, Publisher, Subscriber, ... and write a bunch of test cases to guarantee the library works correctly no matter technology in use.

My reference is from the book https://www.enterpriseintegrationpatterns.com/ and the library https://watermill.io/