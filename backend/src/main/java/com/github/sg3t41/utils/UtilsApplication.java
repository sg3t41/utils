package com.github.sg3t41.utils;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.CommandLineRunner;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.context.annotation.Bean;
import org.springframework.jdbc.core.JdbcTemplate;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

@SpringBootApplication
@RestController
public class UtilsApplication {
    
    @Autowired
    private JdbcTemplate jdbcTemplate;
    
    public static void main(String[] args) {
        SpringApplication.run(UtilsApplication.class, args);
    }
    
    @GetMapping("/")
    public String getDatabaseTime() {
        try {
            return "Database Time: " + jdbcTemplate.queryForObject("SELECT NOW()", String.class);
        } catch (Exception e) {
            return "Database Connection Error: " + e.getMessage();
        }
    }
    
    @Bean
    public CommandLineRunner run(JdbcTemplate jdbcTemplate) {
        return args -> {
            try {
                String result = jdbcTemplate.queryForObject("SELECT NOW()", String.class);
                System.out.println("Database Ping Successful: Current database time is " + result);
            } catch (Exception e) {
                System.err.println("Database Ping Failed: " + e.getMessage());
            }
        };
    }
}
