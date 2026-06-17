import 'dart:convert';
import 'package:http/http.dart' as http;
import '../models/user.dart';

class AuthService {
  // Use 10.0.2.2 for Android Emulator to connect to localhost on the host machine.
  // Use localhost for Web or iOS Simulator.
  // For a real app, this should be moved to a configuration file or .env.
  static const String baseUrl = 'http://localhost:8080';

  Future<Map<String, dynamic>> signup(String email, String password) async {
    final response = await http.post(
      Uri.parse('$baseUrl/auth/signup'),
      headers: {'Content-Type': 'application/json'},
      body: jsonEncode({
        'email': email,
        'password': password,
      }),
    );

    final data = jsonDecode(response.body);

    if (response.statusCode == 201) {
      return {'success': true, 'message': data['message']};
    } else {
      return {'success': false, 'error': data['error'] ?? 'Signup failed'};
    }
  }

  Future<Map<String, dynamic>> login(String email, String password) async {
    final response = await http.post(
      Uri.parse('$baseUrl/auth/login'),
      headers: {'Content-Type': 'application/json'},
      body: jsonEncode({
        'email': email,
        'password': password,
      }),
    );

    final data = jsonDecode(response.body);

    if (response.statusCode == 200) {
      return {
        'success': true,
        'token': data['data']['token'],
        'user': User.fromJson(data['data']['user']),
      };
    } else {
      return {'success': false, 'error': data['error'] ?? 'Login failed'};
    }
  }
}
