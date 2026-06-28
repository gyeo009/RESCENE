package com.example.rescene;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;


import java.time.ZoneOffset;
import java.util.TimeZone;

@SpringBootApplication
public class ResceneApplication {
    public static void main(String[] args) {
        TimeZone.setDefault(TimeZone.getTimeZone(ZoneOffset.UTC));
        SpringApplication.run(ResceneApplication.class, args);
    }
}