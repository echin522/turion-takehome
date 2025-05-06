## Note about the io-processors package
What you see in the io-processors package is a pattern for processing data using simple inputs and outputs. At first glance it might seem cumbersome to implement it in code, but the main advantage of using it as a package is that I can add metrics to the readers, writers, and processor itself very easily. I also personally found myself configuring the same Kafka Consumers and Publishers over and over again, and so having a generic Kafka "Reader" and "Writer" was preferable to me.

