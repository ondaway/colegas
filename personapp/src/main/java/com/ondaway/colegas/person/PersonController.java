package com.ondaway.colegas.person;

import com.rabbitmq.client.AMQP;
import com.rabbitmq.client.Channel;
import com.rabbitmq.client.Connection;
import com.rabbitmq.client.ConnectionFactory;
import com.rabbitmq.client.QueueingConsumer;
import com.rabbitmq.client.AMQP.BasicProperties;
import java.util.UUID;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestMethod;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping("person")
public class PersonController {
    private final Connection       connection;
    private final Channel          channel;
    private final String           requestQueueName = "rpc_queue";
    private final String           replyQueueName;
    private final QueueingConsumer consumer;
    
    
    public PersonController() throws Exception {
        ConnectionFactory factory = new ConnectionFactory();
        factory.setHost("localhost");
        connection = factory.newConnection();
        channel = connection.createChannel();

        replyQueueName = channel.queueDeclare().getQueue();
        consumer = new QueueingConsumer(channel);
        channel.basicConsume(replyQueueName, true, consumer);
    }
    
    
    @RequestMapping(value="/{id}", method=RequestMethod.GET)
    public Person find(@PathVariable String id) throws Exception {
        String ret = personRpc(id);
        return new Person(id, ret, "");
    }

    private String personRpc(String id) throws Exception {
        String correlationId = UUID.randomUUID().toString();
        String response = null;
        
        BasicProperties props = new BasicProperties.Builder()
                .correlationId(correlationId)
                .replyTo(replyQueueName)
                .build();

        channel.basicPublish("", requestQueueName, props, id.getBytes("UTF-8"));
        while (true) {
            QueueingConsumer.Delivery delivery = consumer.nextDelivery();
            if (delivery.getProperties().getCorrelationId().equals(correlationId)) {
                response = new String(delivery.getBody(), "UTF-8");
                break;
            }
        }

        return response;

        
    }
}
