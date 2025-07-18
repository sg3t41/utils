package com.github.sg3t41.utils.healthcheck.usecase;

import com.github.sg3t41.utils.healthcheck.domain.Health;
import lombok.Builder;
import lombok.Getter;

import java.time.LocalDateTime;

/**
 * ユースケース層の出力データ(DTO): ユースケースの結果をインターフェース層に返す
 */
@Getter
@Builder
public class HealthCheckDto {
	private final String status;
	private final LocalDateTime checkedAt;
	private final String details;

	public static HealthCheckDto from(Health health) {
		System.out.println("ユースケース層ログ: HealthCheckDto.from() - HealthドメインオブジェクトをDTOに変換。");
		return HealthCheckDto.builder()
				.status(health.getStatus().name())
				.checkedAt(health.getCheckedAt())
				.details(health.getDetails())
				.build();
	}
}
