package com.github.sg3t41.utils.healthcheck.interfaces;

import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

import com.github.sg3t41.utils.healthcheck.usecase.HealthCheckDto;
import com.github.sg3t41.utils.healthcheck.usecase.HealthCheckUsecase;

import lombok.RequiredArgsConstructor;

/**
 * インターフェース層: 外部からのリクエスト(HTTP)を受け付ける
 */
@RestController
@RequiredArgsConstructor
public class HealthCheckController {

	private final HealthCheckUsecase healthCheckUsecase;

	@GetMapping("/api/health")
	public ResponseEntity<HealthCheckDto> checkHealth() {
		HealthCheckDto outputData = healthCheckUsecase.execute();

		if ("OK".equals(outputData.getStatus())) {
			return ResponseEntity.ok(outputData);
		} else {
			return ResponseEntity.status(503).body(outputData);
		}
	}
}
