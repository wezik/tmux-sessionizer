const std = @import("std");
const json = @import("json");

pub fn newEntry(session_path: []const u8) ConfigEntry {
    return ConfigEntry{
        .session_path = session_path,
        .session_name = session_path,
    };
}

pub const ConfigEntry = struct {
    session_path: []const u8,
    session_name: []const u8,
    windows: []const []const u8 = &[_][]const u8{"shell"},
};

pub const Config = struct {
    entries: []const ConfigEntry = &[_]ConfigEntry{},
};

const configFileName = "config.json";

fn getConfigPath(allocator: std.mem.Allocator) ![]const u8 {
    const env_map = try allocator.create(std.process.EnvMap);
    env_map.* = try std.process.getEnvMap(allocator);

    var value = env_map.get("XDG_CONFIG_HOME") orelse "";
    if (value.len == 0) {
        const home = env_map.get("HOME") orelse ".";
        value = try std.fmt.allocPrint(allocator, "{s}/.config", .{home});
    }

    return try std.fmt.allocPrint(allocator, "{s}/tmux-sessionizer", .{value});
}

fn createDefaultConfig(allocator: std.mem.Allocator) !void {
    var cwd = std.fs.cwd();
    cwd = try cwd.makeOpenPath(try getConfigPath(allocator), .{});
    const f = try cwd.createFile(configFileName, .{});
    const serialized = try json.toPrettySlice(allocator, Config{});
    _ = try f.writer().writeAll(serialized);
}

fn openConfigFile(allocator: std.mem.Allocator) !std.fs.File {
    var cwd = std.fs.cwd();
    cwd = try cwd.makeOpenPath(try getConfigPath(allocator), .{});
    const file = blk: {
        const f = cwd.openFile(configFileName, .{ .mode = .read_write }) catch {
            _ = try createDefaultConfig(allocator);
            break :blk try cwd.openFile(configFileName, .{ .mode = .read_write });
        };
        break :blk f;
    };
    _ = try allocator.dupe(std.fs.File, &[_]std.fs.File{file});
    return file;
}

pub fn fetchConfig(allocator: std.mem.Allocator) !Config {
    const file = try openConfigFile(allocator);
    const file_buffer = try file.readToEndAlloc(allocator, 1024 * 1024);
    const deserialized = try json.fromSlice(allocator, Config, file_buffer);
    return deserialized.value;
}

pub fn saveConfig(allocator: std.mem.Allocator, config: Config) !void {
    const file = try openConfigFile(allocator);
    const serialized = try json.toPrettySlice(allocator, config);
    _ = try file.writer().writeAll(serialized);
}
