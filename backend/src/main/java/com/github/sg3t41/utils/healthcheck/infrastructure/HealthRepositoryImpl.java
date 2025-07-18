package com.github.sg3t41.utils.healthcheck.infrastructure;

import com.github.sg3t41.utils.healthcheck.domain.Health;
import com.github.sg3t41.utils.healthcheck.domain.HealthRepository;
import lombok.RequiredArgsConstructor;
import org.springframework.jdbc.core.JdbcTemplate;
import org.springframework.stereotype.Repository;

import java.time.LocalDateTime;
import java.util.Optional;

/**
 * インフラストラクチャ層: HealthRepositoryの実装
 */
@Repository
@RequiredArgsConstructor
public class HealthRepositoryImpl implements HealthRepository {

	private final JdbcTemplate jdbcTemplate;

	@Override
	public Optional<Health> checkDatabaseHealth() {
		try {
			LocalDateTime dbTime = jdbcTemplate.queryForObject("SELECT NOW()", LocalDateTime.class);
			System.out.println("  > 処理成功: データベース時刻を取得しました (" + dbTime + ")");
			return Optional.of(Health.ok(dbTime));
		} catch (Exception e) {
			System.out.println("  > 処理失敗: データベース接続中にエラーが発生しました。");
			System.err.println(e.getMessage());
			return Optional.of(Health.error(e.getMessage()));
		}
	}
}
