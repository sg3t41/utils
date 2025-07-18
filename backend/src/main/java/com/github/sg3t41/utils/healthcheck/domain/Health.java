package com.github.sg3t41.utils.healthcheck.domain;

import lombok.AllArgsConstructor;
import lombok.Getter;

import java.time.LocalDateTime;

/**
 * ドメインエンティティ: アプリケーションの健康状態を表す
 */
@Getter
@AllArgsConstructor
public class Health {

	private final Status status;
	private final LocalDateTime checkedAt;
	private final String details;

	public enum Status {
		OK,
		ERROR
	}

	/**
	 * データベース接続が正常な場合のHealthオブジェクトを生成するファクトリメソッド
	 */
	public static Health ok(LocalDateTime dbTime) {
		System.out.println("ドメイン層ログ: Health.ok() - 正常状態のHealthオブジェクトを生成。");
		return new Health(Status.OK, dbTime, "データベース接続は正常です。");
	}

	/**
	 * データベース接続に問題がある場合のHealthオブジェクトを生成するファクトリメソッド
	 */
	public static Health error(String errorMessage) {
		System.out.println("ドメイン層ログ: Health.error() - 異常状態のHealthオブジェクトを生成。");
		return new Health(Status.ERROR, LocalDateTime.now(), "データベース接続エラー: " + errorMessage);
	}
}
