use std::error::Error;
use lettre::message::{header, Message};
use lettre::{SmtpTransport, Transport};

use crate::types::event::TaskCreated;

pub trait NotificationService {
    fn send_task_created_notification(&self, event: &TaskCreated) -> Result<(), Box<dyn Error>>;
}

pub struct EmailNotificationService {
    smtp_host: String,
    smtp_port: u16,
}

impl EmailNotificationService {
    pub fn new(smtp_host: String, smtp_port: u16) -> Self {
        EmailNotificationService {
            smtp_host,
            smtp_port,
        }
    }
}

impl NotificationService for EmailNotificationService {
    fn send_task_created_notification(&self, event: &TaskCreated) -> Result<(), Box<dyn Error>> {
        let task = event.task.as_ref().ok_or("task is missing in event")?;
        let subject = format!("Task Created: {}", task.title);
        let body = format!("A new task has been created:\n\nID: {}\nTitle: {}\n", task.id, task.title);
        let email = Message::builder()
            .from("sender@example.com".parse().unwrap())//TODO
            .to("receiver@example.com".parse().unwrap())//TODO
            .subject(subject)
            .header(header::ContentType::TEXT_PLAIN)
            .body(body)?;
        let mailer = SmtpTransport::builder_dangerous(&self.smtp_host)
            .port(self.smtp_port)
            .build();
        mailer.send(&email)?;
        Ok(())
    }
}
