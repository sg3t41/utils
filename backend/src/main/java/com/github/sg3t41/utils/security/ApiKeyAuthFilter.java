package com.github.sg3t41.utils.security;

import java.io.IOException;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.core.authority.AuthorityUtils;
import org.springframework.security.core.context.SecurityContext;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.stereotype.Component;
import org.springframework.util.StringUtils;
import org.springframework.web.filter.OncePerRequestFilter;

import jakarta.servlet.FilterChain;
import jakarta.servlet.ServletException;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;

@Component
public class ApiKeyAuthFilter extends OncePerRequestFilter {

	@Value("${API_SECRET_KEY}")
	private String secretKey;

	private static final String API_KEY_SCHEME = "ApiKey ";

	@Override
	protected void doFilterInternal(HttpServletRequest request, HttpServletResponse response, FilterChain filterChain)
			throws ServletException, IOException {

		String authHeader = request.getHeader("Authorization");

		if (authHeader != null && authHeader.startsWith(API_KEY_SCHEME)) {
			String requestKey = authHeader.substring(API_KEY_SCHEME.length());

			// キーが空でなく、かつ正しいキーであるかを検証
			if (StringUtils.hasText(requestKey) && secretKey.equals(requestKey)) {
				SecurityContext context = SecurityContextHolder.createEmptyContext();
				var authentication = new UsernamePasswordAuthenticationToken("api-user", null,
						AuthorityUtils.createAuthorityList("ROLE_USER"));
				context.setAuthentication(authentication);
				SecurityContextHolder.setContext(context);
			}
		}

		filterChain.doFilter(request, response);
	}
}
