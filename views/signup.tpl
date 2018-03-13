<html>
	<body>
		<h1>signup</h1>
		<form action="/user/signup" method="post">
			用户名:<input type="text" name="username"><br>
			密码:<input type="password" name="password"><br>
			个人介绍:<input type="text" name="introduction"><br>
			<input type="submit" value="注册">
		</form>
		<div style="position:absolute;bottom:0">{{.content}}<div>
	</body>
</html>