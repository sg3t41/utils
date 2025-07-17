package com.github.sg3t41.utils.controller;

import org.springframework.jdbc.core.JdbcTemplate;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

import lombok.RequiredArgsConstructor;

@RestController
@RequiredArgsConstructor
public class HealthCheckController {

	private final JdbcTemplate jdbcTemplate;

	@GetMapping("/api/health")
	public String getDatabaseTime() {
		try {
			return "Database Time: " + jdbcTemplate.queryForObject("SELECT NOW()", String.class);
		} catch (Exception e) {
			return "Database Connection Error: " + e.getMessage();
		}
	}
}
