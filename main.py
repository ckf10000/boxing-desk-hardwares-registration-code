import os
import subprocess
import hashlib
import random
import time
from typing import List

class HardwareInfo:
    def __init__(self):
        self.int_code = [0] * 127
        self.int_number = [0] * 25
        self.charcode = [0] * 25  # Using 0 as byte equivalent in Python

    def get_cpu(self) -> str:
        try:
            result = subprocess.run(
                ["powershell", "Get-WmiObject Win32_Processor | Select-Object -ExpandProperty ProcessorId"],
                stdout=subprocess.PIPE,
                stderr=subprocess.PIPE,
                check=True,
                universal_newlines=True  # 替换text=True
            )
            return result.stdout.strip()
        except subprocess.SubprocessError as e:
            print(f"获取CPU ID失败: {e}")
            return "BFEBFBFF000"

    def get_disk_volume_serial_number(self) -> str:
        disk_device_id = "C:"
        power_shell_cmd = (
            f"Get-WmiObject Win32_LogicalDisk -Filter \"DeviceID='{disk_device_id}'\" | "
            "Select-Object DeviceID, VolumeSerialNumber"
        )
        try:
            result = subprocess.run(
                ["powershell", power_shell_cmd],
                stdout=subprocess.PIPE,
                stderr=subprocess.PIPE,
                check=True,
                universal_newlines=True  # 替换text=True
            )
            output = result.stdout
            volume_serial_number_slice = output.split(disk_device_id)
            if len(volume_serial_number_slice) != 2:
                print(f"获取 C盘的卷序列号（VolumeSerialNumber）:\n{output}")
                return "00000000"
            return volume_serial_number_slice[1].strip()
        except subprocess.SubprocessError as e:
            print(f"获取 C盘的卷序列号（VolumeSerialNumber）失败: {e}")
            return "00000000"

    def get_m_num(self) -> str:
        cpu = self.get_cpu()
        disk = self.get_disk_volume_serial_number()

        if cpu == "BFEBFBFF000" or disk == "00000000":
            random.seed(time.time_ns())
            suffix = f"{random.randint(0, 0xFFFFFF):06X}"
            return suffix + "0" * 18

        combined = cpu + disk
        if len(combined) > 24:
            return combined[:24]
        return combined

    def set_int_code(self) -> None:
        for i in range(1, len(self.int_code)):
            self.int_code[i] = i % 9

    def get_r_num(self) -> str:
        self.set_int_code()
        m_num = self.get_m_num()

        # Fill charcode array (index starts at 1)
        for i in range(1, len(self.charcode)):
            if i - 1 < len(m_num):
                self.charcode[i] = ord(m_num[i - 1])
            else:
                self.charcode[i] = ord('0')

        # Calculate int_number array (index starts at 1)
        for j in range(1, len(self.int_number)):
            char_code = self.charcode[j]
            self.int_number[j] = self.int_code[char_code] + char_code

        # Build ASCII string
        ascii_str = ""
        for k in range(1, len(self.int_number)):
            code = self.int_number[k]
            if 48 <= code <= 57:  # 0-9
                ascii_str += chr(code)
            elif 65 <= code <= 90:  # A-Z
                ascii_str += chr(code)
            elif 97 <= code <= 122:  # a-z
                ascii_str += chr(code)
            elif code <= 122:
                ascii_str += chr(code - 9)
            else:
                ascii_str += chr(code - 10)

        # SHA1 hash calculation (first 10 chars uppercase)
        return self.generate_sha1(ascii_str)

    @staticmethod
    def generate_sha1(input_str: str) -> str:
        sha1_hash = hashlib.sha1(input_str.encode()).hexdigest()
        return sha1_hash[:10].upper()


def main():
    hardware = HardwareInfo()

    # Generate and print registration code
    r_num = hardware.get_r_num()
    print(f"注册码: {r_num}")

    # More professional implementation
    print("\n按任意键退出...")
    input()  # Wait for any key press


if __name__ == "__main__":
    main()