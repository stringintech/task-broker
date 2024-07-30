use prost::Message;
use lapin::{options::*, types::FieldTable, Connection, ConnectionProperties};
use futures_util::stream::StreamExt;
use std::error::Error;
use std::io::Cursor;

mod message {
    include!(concat!(env!("OUT_DIR"), "/message.rs"));
}

use message::TaskMessage;

#[tokio::main]
async fn main() -> Result<(), Box<dyn Error>> {
    let addr = "amqp://guest:guest@localhost:5672//";
    let conn = Connection::connect(&addr, ConnectionProperties::default()).await?;
    let channel = conn.create_channel().await?;

    let _queue = channel.queue_declare(
        "hello",
        QueueDeclareOptions::default(),
        FieldTable::default(),
    ).await?;

    let mut consumer = channel.basic_consume(
        "hello",
        "my_consumer",
        BasicConsumeOptions::default(),
        FieldTable::default(),
    ).await?;

    println!("waiting for messages");

    while let Some(delivery) = consumer.next().await {
        match delivery {
            Ok(delivery) => {
                let msg = TaskMessage::decode(&mut Cursor::new(&delivery.data))?;
                println!("received message: {}", msg.content);
                delivery.ack(BasicAckOptions::default()).await?;
            }
            Err(error) => {
                eprintln!("error receiving message: {:?}", error);
            }
        }
    }

    Ok(())
}
