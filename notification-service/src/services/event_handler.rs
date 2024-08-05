use crate::types::event::TaskCreated;
use std::error::Error;
use super::notification_service::NotificationService;

pub struct EventHandler {
    notification_service: Box<dyn NotificationService>,
}

impl EventHandler {
    pub fn new(notification_service: Box<dyn NotificationService>) -> Self {
        EventHandler { notification_service }
    }

    pub fn handle_event_task_created(&self, event: &TaskCreated) -> Result<(), Box<dyn Error>> {
        self.notification_service.send_task_created_notification(event)
    }
}