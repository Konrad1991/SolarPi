use chrono::Local;
use std::fs::{self, OpenOptions};
use std::io::{self, Write};

pub fn log(message: &str) {
    eprintln!("{}", message);
}

pub fn log_file(message: &str, log_file_path: &str) {
    eprintln!("{}", message);
    if let Ok(mut file) = OpenOptions::new().create(true).append(true).open(log_file_path) {
        if let Err(err) = writeln!(file, "[{}] {}", Local::now(), message) {
            eprintln!("Error writing to log file: {}", err);
        }
    } else {
        eprintln!("Failed to open or create log file: {}", log_file_path);
    }
}

pub enum LogLevel {
    Info,
    Warning,
    Error,
}

pub fn log_level(level: LogLevel, message: &str) {
    match level {
        LogLevel::Info => println!("[INFO] {}", message),
        LogLevel::Warning => println!("[WARNING] {}", message),
        LogLevel::Error => println!("[ERROR] {}", message),
    }
}
