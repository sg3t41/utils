package com.github.sg3t41.utils.healthcheck.usecase;

import com.github.sg3t41.utils.healthcheck.domain.Health;
import com.github.sg3t41.utils.healthcheck.domain.HealthRepository;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

/**
 * ユースケース層: ヘルスチェックのビジネスフローを実装する
 *
 * - アプリケーション固有のロジックをここに記述する。
 * - ドメイン層のオブジェクトやリポジトリを使ってビジネスフローを実現する。
 * - Springの@Serviceアノテーションを付与して、DIコンテナに登録する。
 */
@Service
@RequiredArgsConstructor
public class HealthCheckUsecase {

	private final HealthRepository healthRepository;

	/**
	 * ヘルスチェックを実行し、結果をDTOで返す
	 * 
	 * @return ヘルスチェック結果のDTO
	 */
	@Transactional(readOnly = true) // 読み取り専用のトランザクション
	public HealthCheckDto execute() {
		System.out.println("ユースケース層ログ: HealthCheckUsecase.execute() - 呼び出し。");
		System.out.println("  > 本来の処理: ドメインオブジェクトやリポジトリを使い、一連のビジネスフローを実行する。");

		// 1. リポジトリを呼び出して、インフラ層からデータを取得する
		Health health = healthRepository.checkDatabaseHealth()
				.orElse(Health.error("リポジトリが空の結果を返しました。"));

		// 2. (もしあれば)ドメインロジックを実行する
		// 例: health.someDomainLogic();

		// 3. 結果をDTOに変換して、インターフェース層に返す
		return HealthCheckDto.from(health);
	}
}
