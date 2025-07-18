package com.github.sg3t41.utils.healthcheck.domain;

import java.util.Optional;

/**
 * ドメインリポジトリのインターフェース: Healthエンティティの永続化を抽象化する
 */
public interface HealthRepository {

	/**
	 * 現在のデータベースの状態（時刻）を取得する
	 */
	Optional<Health> checkDatabaseHealth();
}
