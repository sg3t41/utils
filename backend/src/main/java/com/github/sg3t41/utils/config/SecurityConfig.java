package com.github.sg3t41.utils.config;

import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.config.http.SessionCreationPolicy;
import org.springframework.security.web.SecurityFilterChain;
import org.springframework.security.web.authentication.UsernamePasswordAuthenticationFilter;

import com.github.sg3t41.utils.security.ApiKeyAuthFilter;

import lombok.RequiredArgsConstructor;

@Configuration
@RequiredArgsConstructor
public class SecurityConfig {

	private final ApiKeyAuthFilter apiKeyAuthFilter;

	@Bean
	public SecurityFilterChain securityFilterChain(HttpSecurity http) throws Exception {
		http
				// Spring Securityのデフォルト認証とCORS設定を無効化
				.csrf(csrf -> csrf.disable())
				.httpBasic(basic -> basic.disable())
				.formLogin(form -> form.disable())
				.logout(logout -> logout.disable())
				// REST APIなのでセッションはステートレスにする
				.sessionManagement(session -> session.sessionCreationPolicy(SessionCreationPolicy.STATELESS))
				// URLごとに認証を設定
				.authorizeHttpRequests(auth -> auth
						.requestMatchers("/api/health").permitAll() // ヘルスチェック用エンドポイントは認証不要
						.requestMatchers("/api/**").authenticated() // それ以外の/api/**は認証が必要
						.anyRequest().permitAll() // その他のリクエストはすべて許可
				)
				// 自作のAPIキー認証フィルターを追加
				.addFilterBefore(apiKeyAuthFilter, UsernamePasswordAuthenticationFilter.class);

		return http.build();
	}
}
