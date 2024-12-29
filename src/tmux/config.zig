const std = @import("std");
const json = @import("json");
const domain = @import("domain.zig");

const configFileName = "config.json";

fn getConfigPath(allocator: std.mem.Allocator) ![]const u8 {
    const env_map = try allocator.create(std.process.EnvMap);
    env_map.* = try std.process.getEnvMap(allocator);

    var value = env_map.get("XDG_CONFIG_HOME") orelse "";
    if (value.len == 0) {
        const home = env_map.get("HOME") orelse ".";
        value = try std.fmt.allocPrint(allocator, "{s}/.config", .{home});
    }

    const result = try std.fmt.allocPrint(allocator, "{s}/tmux-sessionizer", .{value});
    return result;
}

pub fn getFullConfigPath(allocator: std.mem.Allocator) ![]const u8 {
    return try std.fmt.allocPrint(allocator, "{s}/{s}", .{ try getConfigPath(allocator), configFileName });
}

fn createDefaultConfig(allocator: std.mem.Allocator) !void {
    var cwd = std.fs.cwd();
    cwd = try cwd.makeOpenPath(try getConfigPath(allocator), .{});
    const f = try cwd.createFile(configFileName, .{});
    const serialized = try json.toPrettySlice(allocator, [_]domain.TmuxSession{});
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

pub fn fetchConfig(allocator: std.mem.Allocator) ![]const domain.TmuxSession {
    const file = try openConfigFile(allocator);
    const file_buffer = try file.readToEndAlloc(allocator, 1024 * 1024);
    const deserialized = try json.fromSlice(allocator, []const domain.TmuxSession, file_buffer);
    return deserialized.value;
}

pub fn saveConfig(allocator: std.mem.Allocator, config: []domain.TmuxSession) !void {
    // TODO: okay this is bad but it's a quick fix for now
    // * when removing a session, the config file needs to be truncated however
    // * as far as I've seen there are some weird inconsistencies with how zig behaves,
    // * need to explore this further.
    //
    // * also, maybe it's just better to create a new file and overwrite it every time
    _ = try createDefaultConfig(allocator);

    const file = try openConfigFile(allocator);
    const serialized = try json.toPrettySlice(allocator, config);
    _ = try file.writer().writeAll(serialized);
}
